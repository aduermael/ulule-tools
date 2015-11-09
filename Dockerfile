FROM golang:1.5.1

# Go package experimental vendoring
ENV GO15VENDOREXPERIMENT 1

ADD ./go /go
CMD ["/bin/bash"]
