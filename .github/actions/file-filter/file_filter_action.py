#!/usr/bin/env python3

import argparse
import fnmatch
import json
import os
import sys

import github


def parse_patterns(patterns_input: str) -> list[str]:
    """Parse glob patterns from newline- or space-separated string input."""
    if not patterns_input:
        return []

    patterns = [p.strip() for p in patterns_input.split() if p.strip()]

    if not patterns:
        raise ValueError('No valid patterns found in input')

    return patterns


def get_changed_files(
    github_client: github.Github, repo_name: str,
) -> list[str]:
    """Get list of changed files between base and head refs."""
    try:
        repo = github_client.get_repo(repo_name)

        # Try to get PR context first
        if 'GITHUB_EVENT_PATH' in os.environ:
            try:
                with open(os.environ['GITHUB_EVENT_PATH']) as f:
                    event_data = json.load(f)

                if 'pull_request' in event_data:
                    pr_number = event_data['pull_request']['number']
                    pr = repo.get_pull(pr_number)
                    files = pr.get_files()
                    return [f.filename for f in files]
            except (FileNotFoundError, KeyError, json.JSONDecodeError):
                pass

        base_ref = os.environ.get('GITHUB_BASE_REF', 'main')
        head_ref = os.environ.get('GITHUB_HEAD_REF', os.environ.get('GITHUB_SHA'))

        if head_ref:
            comparison = repo.compare(base_ref, head_ref)
            return [f.filename for f in comparison.files]

        print('Warning: Could not determine changed files', file=sys.stderr)
        return []
    except github.GithubException as e:
        print(f'GitHub API error: {e}', file=sys.stderr)
        return []
    except Exception as e:
        print(f'Error getting changed files: {e}', file=sys.stderr)
        return []


def match_files(files: list[str], patterns: list[str], exclude: bool) -> list[str]:
    """Match files against glob patterns."""
    matches = []

    for file_path in files:
        for pattern in patterns:
            if (fnmatch.fnmatch(file_path, pattern) and not exclude) or exclude:
                matches.append(file_path)
                break

    return matches


def set_output(name: str, value: str) -> None:
    """Set GitHub Actions output."""
    if 'GITHUB_OUTPUT' in os.environ:
        with open(os.environ['GITHUB_OUTPUT'], 'a') as f:
            f.write(f'{name}={value}\n')
    else:
        # Fallback for older runners
        print(f'::set-output name={name}::{value}')


def main() -> None:
    """Main function."""
    parser = argparse.ArgumentParser(
        description=(
            'A utility script to retrieve the list of changed files in the '
            'current PR, intended to be run as part of a GitHub Actions '
            'pipeline.'
        ),
    )
    parser.parse_args()

    try:
        patterns_input = os.environ.get('INPUT_PATTERNS', '')
        token = os.environ.get('INPUT_TOKEN', os.environ.get('GITHUB_TOKEN', ''))
        exclude = os.environ.get('INPUT_EXCLUDE')
        repo_name = os.environ.get('GITHUB_REPOSITORY', '')
        event_type = os.environ.get('GITHUB_EVENT_NAME')

        print(f'Event type: {event_type}')
        if event_type in ('schedule',):
            print(f'Skipping file check for event_type={event_type}')
            set_output('matches', 'true')
            set_output('count', '')
            set_output('files', '')
            sys.exit(0)

        if not patterns_input:
            print('Error: No patterns provided', file=sys.stderr)
            sys.exit(1)

        if not token:
            print('Error: No GitHub token provided', file=sys.stderr)
            sys.exit(1)

        if not repo_name:
            print('Error: No repository name found', file=sys.stderr)
            sys.exit(1)

        if exclude and exclude not in ('true', 'false'):
            print('Error: exclude must be one of: true, false', file=sys.stderr)
            sys.exit(1)

        patterns = parse_patterns(patterns_input)
        print(f'Parsed patterns: {patterns}')

        github_client = github.Github(token)
        changed_files = get_changed_files(github_client, repo_name)
        matched_files = match_files(changed_files, patterns, exclude == 'true')

        print(f'Has matches? {"true" if matched_files else "false"}')
        print(f'Matched files: {matched_files}')

        set_output('matches', 'true' if matched_files else 'false')
        set_output('count', str(len(matched_files)))
        set_output('files', json.dumps(matched_files))
        sys.exit(0)
    except Exception as e:
        print(f'Error: {e}', file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()
