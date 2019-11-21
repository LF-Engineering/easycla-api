FROM python:3.7.4
RUN mkdir cla logs

ADD requirements.txt *.py /
RUN pip install --no-cache-dir -r requirements.txt


ADD cla /cla 

CMD [ "python", "data_migration.py" ]