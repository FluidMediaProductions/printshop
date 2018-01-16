FROM scratch

COPY printshop /
COPY config/ /config
COPY imgs/ /imgs

CMD ["/printshop"]
