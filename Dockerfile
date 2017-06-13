FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./d.ims.io /
CMD ["/d.ims.io"]
