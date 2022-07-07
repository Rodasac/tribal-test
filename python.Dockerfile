FROM python:alpine
ENV FLASK_APP main
ENV FLASK_ENV development

WORKDIR /usr/src/app

COPY python/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY python .

EXPOSE 5000

CMD [ "flask", "run", "--host=0.0.0.0" ]