
FROM python:2.7.9-wheezy

WORKDIR /opt

COPY ./Torshammer1.0/torshammer.py .

RUN pip install PySocks
RUN pip install Terminal

CMD ["python2", "/opt/torshammer.py", "-t localhost", "-p 7000", "-r 30"] 
