FROM tcloud.hub/base/python:3.8
RUN apt-get update && apt-get install vim -y && apt-get install sshpass
RUN mkdir -p /home/py_xops && mkdir -p ~/.pip/
COPY ./.pip/pip.conf ~/.pip/
WORKDIR /home/py_xops
COPY . /home/py_xops
RUN pip install -r /home/py_xops/requirements.txt
EXPOSE 8000
CMD ["python", "/home/py_xops/main.py", "-c", "/home/py_xops/api.cfg"]
