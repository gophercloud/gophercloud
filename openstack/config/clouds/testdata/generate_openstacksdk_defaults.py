"""Generate configuration expectations with openstacksdk (optional test tooling).

Install the versions used by the checked-in fixture in a virtual environment:
    pip install openstacksdk==4.20.0 os-service-types==1.9.0 PyYAML==6.0.3
Then regenerate:
    python generate_openstacksdk_defaults.py > openstacksdk-defaults.yaml
Go tests consume the checked-in YAML; they do not require Python or an SDK install.
"""
import importlib.metadata
import tempfile
from pathlib import Path

import yaml
from openstack.config.loader import OpenStackConfig

defaults = {'region': 'mars', 'service': 'block-storage', 'cloud': {'auth': {'auth_url': 'https://example.org/v3'}, 'region_name': 'mars'}}
cases = []
def add(name, cloud, service='block-storage', secure=None, public=None, region='mars'):
    case = dict(name=name, cloud=cloud)
    for key, value, default in [('service', service, defaults['service']), ('region', region, defaults['region']), ('secure', secure, None), ('public', public, None)]:
        if value != default:
            case[key] = value
    cases.append(case)

aliases = ['block_storage', 'volumev3', 'volumev2', 'volume', 'block_store']
for alias in aliases:
    add('alias-' + alias, {alias + '_default_microversion': '3.60'})
for i, alias in enumerate(aliases[:-1]):
    add('precedence-' + alias, {key + '_default_microversion': '3.' + str(60 + j) for j, key in enumerate(aliases[i:])}, service='volumev3')
add('empty-canonical', {'block_storage_default_microversion':'', 'volumev3_default_microversion':'3.60'})
add('null-canonical', {'block_storage_default_microversion':None, 'volumev3_default_microversion':'3.60'})
add('missing', {'compute_default_microversion':'2.87'})
add('global-is-not-default', {'default_microversion':'3.60'})
add('unknown-service', {'future_service_default_microversion':'1.4'}, service='future-service')
add('hyphenated-key', {'block-storage-default-microversion':'3.60'}, service='block_storage')
add('quoted-trailing-zero', {'compute_default_microversion':'2.10'}, service='compute')
add('secure-override', {'compute_default_microversion':'2.87'}, service='compute', secure={'compute_default_microversion':'2.79'})
add('secure-null', {'block_storage_default_microversion':'3.60','volumev3_default_microversion':'3.50'}, secure={'block_storage_default_microversion':None})
add('regional-override', {'compute_default_microversion':'2.87','regions':[{'name':'mars','values':{'compute_default_microversion':'2.79'}}]}, service='compute')
add('regional-null-inherits', {'compute_default_microversion':'2.87','regions':[{'name':'mars','values':{'compute_default_microversion':None}}]}, service='compute')
add('regional-empty-clears', {'compute_default_microversion':'2.87','regions':[{'name':'mars','values':{'compute_default_microversion':''}}]}, service='compute')
add('regional-alias-vs-canonical', {'block_storage_default_microversion':'3.60','regions':[{'name':'mars','values':{'volumev3_default_microversion':'3.50'}}]})
add('secure-and-region', {'compute_default_microversion':'2.87','regions':[{'name':'mars','values':{'compute_default_microversion':'2.79'}}]},service='compute',secure={'compute_default_microversion':'2.80'})
add('profile-default', {'profile':'example'},service='compute',public={'compute_default_microversion':'2.87'})
add('profile-cloud-override', {'profile':'example','compute_default_microversion':'2.79'},service='compute',public={'compute_default_microversion':'2.87'})
add('profile-secure-override', {'profile':'example','compute_default_microversion':'2.79'},service='compute',public={'compute_default_microversion':'2.87'},secure={'compute_default_microversion':'2.80'})

with tempfile.TemporaryDirectory() as directory:
    root = Path(directory)
    for case in cases:
        cloud = {**defaults['cloud'], **case['cloud']}
        (root/'clouds.yaml').write_text(yaml.safe_dump({'clouds': {'test': cloud}}))
        (root/'secure.yaml').write_text(yaml.safe_dump({'clouds': {'test': case.get('secure', {})}}))
        (root/'clouds-public.yaml').write_text(yaml.safe_dump({'public-clouds': {'example': case.get('public', {})}}))
        config = OpenStackConfig(config_files=[str(root/'clouds.yaml')], secure_files=[str(root/'secure.yaml')],vendor_files=[str(root/'clouds-public.yaml')],load_envvars=False)
        case['expected'] = config.get_one('test', region_name=case.get('region', defaults['region']), validate=False).get_default_microversion(case.get('service', defaults['service'])) or ''

print(yaml.safe_dump({'openstacksdk':importlib.metadata.version('openstacksdk'),'os_service_types':importlib.metadata.version('os-service-types'),'defaults':defaults,'cases':cases},sort_keys=False),end='')
