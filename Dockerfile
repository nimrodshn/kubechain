FROM golang:latest

RUN mkdir -p $GOPATH/src/github.com/nimrodshn/kubechain

WORKDIR $GOPATH/src/github.com/nimrodshn/kubechain

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh 
    
COPY . $GOPATH/src/github.com/nimrodshn/kubechain

RUN dep ensure -v

RUN make

CMD ["./kubechain"]

