FROM ubuntu:latest
MAINTAINER Xavier Pivan <xavier.pivan.ds@gmail.com>

ENV HOME /root
ENV HDF5_DIR /usr/local/hdf5

LABEL "eudat.gef.service.name"="enes-usecase"
LABEL "eudat.gef.service.description"="Download and average all the dataset together" 
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory" 
LABEL "eudat.gef.service.input.1.path"="/root/input"
LABEL "eudat.gef.service.input.1.type"="string"
LABEL "eudat.gef.service.input.1.filename"="list_url.txt"
LABEL "eudat.gef.service.output.1.name"="Output Directory" 
LABEL "eudat.gef.service.output.1.path"="/root/output"

RUN apt-get update && apt-get -yq install gcc \
                                python2.7 \
                                python-pip \
                                python-dev \
								build-essential \
								wget \
								libhdf5-serial-dev \
								tar \
								bzip2 \
								m4 \
                                git \
								zlib1g-dev \
								libopenmpi-dev

RUN pip install --upgrade pip 

COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt

COPY install_netCDF4_python.sh install_netCDF4_python.sh

RUN chmod +x install_netCDF4_python.sh && \
    ./install_netCDF4_python.sh

RUN git clone https://github.com/cerfacs-globc/icclim && \
    cd icclim && \
    gcc -fPIC -g -c -Wall ./icclim/libC.c -o ./icclim/libC.o && \
    gcc -shared -o ./icclim/libC.so ./icclim/libC.o && \
    python setup.py install --user

COPY ENES_UseCase_icclim.py ENES_UseCase_icclim.py

CMD python ENES_UseCase_icclim.py