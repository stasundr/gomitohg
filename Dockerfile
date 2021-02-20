FROM neilotoole/xcgo:latest

WORKDIR /mitohg

COPY . .

RUN git clone https://github.com/smarco/WFA
RUN cd WFA && make clean all

RUN apt-get install -y libjson-c-dev

RUN cd /mitohg/bridge;\
  gcc -O2 -I../WFA -c wfa_bridge.c;\
  ar q libwfabridge.a wfa_bridge.o;\
  cd /mitohg && go build -o mitohg