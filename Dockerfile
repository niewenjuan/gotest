FROM 100.125.0.198:20202/hwcse/as-go:1.8.5

COPY ./test0518go /home
COPY ./conf /home/conf
RUN chmod +x /home/test0518go

CMD ["/home/test0518go"]