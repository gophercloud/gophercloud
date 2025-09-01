FROM python:3.11-slim
WORKDIR /app
RUN pip install 'PyGithub~=2.7.0' 'requests~=2.32'
COPY file_filter_action.py ./
ENTRYPOINT ["python", "/app/file_filter_action.py"]
