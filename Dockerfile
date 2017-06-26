FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./static /static
ADD ./d.ims.io /
CMD ["/d.ims.io"]
