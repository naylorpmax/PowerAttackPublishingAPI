FROM python:3.9

LABEL author="Max Naylor"

WORKDIR /usr/src/app

COPY requirements.txt ./

RUN python -m pip install --upgrade pip && \
	pip install --no-cache-dir -r requirements.txt
