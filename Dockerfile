FROM golang:1.11 as go_builder

COPY . /go/src/malice-new
WORKDIR /go/src/malice-new
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan


FROM debian:jessie

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/avg.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"


RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware


# Install Requirements
RUN buildDeps='ca-certificates curl' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 lib32z1 --no-install-recommends \
  && echo "===> Install AVG..." \
  && curl -Ls http://download.avgfree.com/filedir/inst/avg2013flx-r3118-a6926.i386.deb > /tmp/avg.deb \
  && dpkg -i /tmp/avg.deb \
  && /etc/init.d/avgd restart \
  && avgcfgctl -w UpdateVir.sched.Task.Disabled=true \
  && avgcfgctl -w Default.setup.daemonize=false \
  && avgcfgctl -w Default.setup.features.antispam=false \
  && avgcfgctl -w Default.setup.features.oad=false \
  && avgcfgctl -w Default.setup.features.scheduler=false \
  && avgcfgctl -w Default.setup.features.tcpd=false \
  && sed -i 's/Severity=INFO/Severity=None/g' /opt/avg/av/cfg/scand.ini \
  && sed -i 's/Severity=INFO/Severity=None/g' /opt/avg/av/cfg/tcpd.ini \
  && sed -i 's/Severity=INFO/Severity=None/g' /opt/avg/av/cfg/wd.ini 


ARG ZONE_KEY
ENV ZONE_KEY=$ZONE_KEY

ENV ZONE 1.3.0

RUN buildDeps='ca-certificates wget build-essential' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 \
  && echo "===> Install Zoner AV..." \
  # && wget -q -P /tmp http://update.zonerantivirus.com/download/zav-${ZONE}-ubuntu-amd64.deb \
  && wget --progress=bar:force -P /tmp https://github.com/maliceio/malice-av/raw/master/zoner/zav-1.3.0-debian-amd64.deb \
  && dpkg -i /tmp/zav-${ZONE}-debian-amd64.deb; \
  if [ "x$ZONE_KEY" != "x" ]; then \
  echo "===> Updating License Key..."; \
  sed -i "s/UPDATE_KEY.*/UPDATE_KEY = ${ZONE_KEY}/g" /etc/zav/zavd.conf; \
  fi 


RUN mkdir -p /opt/malice
RUN if [ "x$ZONE_KEY" != "x" ]; then \
  echo "===> Update zoner definitions..."; \
  /etc/init.d/zavd update; \
  fi


RUN buildDeps='ca-certificates wget' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 --no-install-recommends \
  && set -x \
  && echo "===> Install F-PROT..." \
  && wget https://github.com/maliceio/malice-av/raw/master/fprot/fp-Linux.x86.32-ws.tar.gz \
  -O /tmp/fp-Linux.x86.32-ws.tar.gz \
  && tar -C /opt -zxvf /tmp/fp-Linux.x86.32-ws.tar.gz \
  && ln -fs /opt/f-prot/fpscan /usr/local/bin/fpscan \
  && ln -fs /opt/f-prot/fpscand /usr/local/sbin/fpscand \
  && ln -fs /opt/f-prot/fpmon /usr/local/sbin/fpmon \
  && cp /opt/f-prot/f-prot.conf.default /opt/f-prot/f-prot.conf \
  && ln -fs /opt/f-prot/f-prot.conf /etc/f-prot.conf \
  && chmod a+x /opt/f-prot/fpscan \
  && chmod u+x /opt/f-prot/fpupdate \
  && ln -fs /opt/f-prot/man_pages/scan-mail.pl.8 /usr/share/man/man8/ 





#install clamAv
RUN apk --update add --no-cache clamav ca-certificates
RUN apk --update add --no-cache -t .build-deps \
  build-base \
  mercurial \
  musl-dev \
  openssl \
  bash \
  wget \
  git 




RUN echo "===> Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg

# Ensure ca-certificates is installed for elasticsearch to use https
RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Update AVG Definitions
RUN mkdir -p /opt/malice && /etc/init.d/avgd restart && avgupdate

# Update F-PROT Definitions
RUN mkdir -p /opt/malice && /opt/f-prot/fpupdate






# # Add EICAR Test Virus File to malware folder
# ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

COPY --from=go_builder /bin/avscan /bin/avscan

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]
CMD ["--help"]