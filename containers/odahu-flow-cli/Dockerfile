#
#    Copyright 2019 EPAM Systems
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#

FROM ubuntu:18.04

ARG MINICONDA_URL=https://repo.anaconda.com/miniconda/Miniconda3-4.5.11-Linux-x86_64.sh

ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
ENV PATH /opt/conda/bin:$PATH

# Install native dependencies
RUN apt-get update --fix-missing && \
    apt-get install -y wget bzip2 ca-certificates curl git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install conda
RUN wget --quiet ${MINICONDA_URL}  -O ~/miniconda.sh && \
    /bin/bash ~/miniconda.sh -b -p /opt/conda && \
    rm ~/miniconda.sh && \
    /opt/conda/bin/conda clean -tipsy && \
    ln -s /opt/conda/etc/profile.d/conda.sh /etc/profile.d/conda.sh && \
    echo ". /opt/conda/etc/profile.d/conda.sh" >> ~/.bashrc && \
    echo "conda activate base" >> ~/.bashrc

RUN pip install --upgrade pip

ARG WORKSPACE=/opt/workspace

COPY packages/ ${WORKSPACE}/

# Install python dependecies
RUN  cd ${WORKSPACE}/sdk/ && pip install -r requirements.txt
RUN  cd ${WORKSPACE}/cli/ && pip install -r requirements.txt

# Install packages
RUN cd ${WORKSPACE}/sdk/ && pip install --no-deps -e .
RUN cd ${WORKSPACE}/cli/ && pip install --no-deps -e .

CMD [ "/bin/bash" ]
