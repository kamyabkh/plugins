FROM golang:1.16 as go_builder

COPY . /go/src/malice-new
WORKDIR /go/src/malice-new
#RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
#RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan
COPY avscan /bin/



# FROM ubuntu:precise


# RUN buildDeps='ca-certificates \
#   build-essential \
#   gdebi-core \
#   libssl-dev \
#   mercurial \
#   git-core \
#   wget' \
#   && apt-get update -qq \
#   && apt-get install -yq $buildDeps \
#   && echo "===> Install Comodo..." \
#   && cd /tmp \
#   #&& wget https://cdn.download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb \
#   && wget http://download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb \
#   && DEBIAN_FRONTEND=noninteractive gdebi -n cav-linux_x64.deb \
#   && DEBIAN_FRONTEND=noninteractive /opt/COMODO/post_setup.sh 




# ADD http://download.comodo.com/av/updates58/sigs/bases/bases.cav /opt/COMODO/scanners/bases.cav


FROM ubuntu:bionic

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/fsecure.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware


ENV FSECURE_VERSION 11.10.68

# Install Requirements
RUN buildDeps='wget rpm ca-certificates' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps lib32stdc++6 psmisc \
  && echo "===> Install F-Secure..." \
  && cd /tmp \
  && wget -q https://download.f-secure.com/corpro/ls/trial/fsls-${FSECURE_VERSION}-rtm.tar.gz \
  && tar zxvf fsls-${FSECURE_VERSION}-rtm.tar.gz \
  && cd fsls-${FSECURE_VERSION}-rtm \
  && chmod a+x fsls-${FSECURE_VERSION} \
  && ./fsls-${FSECURE_VERSION} --auto standalone lang=en --command-line-only \
  && fsav --version \
  && echo "===> Update F-Secure..." \
  && cd /tmp \
  && wget -q http://download.f-secure.com/latest/fsdbupdate9.run \
  && mv fsdbupdate9.run /opt/f-secure/ 
  # && echo "===> Clean up unnecessary files..." \
  # && apt-get purge -y --auto-remove $buildDeps && apt-get clean \
  # && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg

# # Ensure ca-certificates is installed for elasticsearch to use https
# RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates \
#   && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Update F-Secure
RUN echo "===> Update F-Secure Database..." \
  && mkdir -p /opt/malice \
  && /etc/init.d/fsaua start \
  && /etc/init.d/fsupdate start \
  && /opt/f-secure/fsav/bin/dbupdate /opt/f-secure/fsdbupdate9.run; exit 0



# Install Requirements
RUN buildDeps='wget ca-certificates' \
  && DEBIAN_FRONTEND=noninteractive apt-get update -qq \
  && apt-get install -yq $buildDeps 
RUN  echo "===> Install Sophos..." 
# COPY sav-linux-free-9.tgz /tmp/
RUN cd /tmp \
   && wget -q https://github.com/maliceio/malice-av/raw/master/sophos/sav-linux-free-9.tgz \
  && tar xzvf sav-linux-free-9.tgz \
  && ./sophos-av/install.sh /opt/sophos --update-free  --acceptlicence --autostart=False --enableOnBoot=False --automatic --ignore-existing-installation --update-source-type=s \
  && echo "===> Update Sophos..." \
  && mkdir -p /opt/malice \
  && /opt/sophos/update/savupdate.sh \
  && /opt/sophos/bin/savconfig set DisableFeedback true 
  

# https://download.eset.com/com/eset/apps/business/efs/linux/latest/efs.x86_64.bin

# RUN wget https://download.eset.com/com/eset/apps/business/efs/linux/latest/efs.x86_64.bin > /tmp/efs.x86_64.bin\
# && sh ./efs.x86_64.bin\
# && /opt/eset/efs/sbin/lic -k  MCU2-XFVE-CM6S-C2MK-T56D




#/opt/eset/efs/lib/egui -v  version eset
#/opt/eset/efs/bin/upd --update --server=192.168.1.2:2221
# /opt/eset/efs/sbin/lic -k XXXX-XXXX-XXXX-XXXX-XXXX


FROM debian:jessie

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/avg.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"



# # Install Requirements
# RUN buildDeps='ca-certificates curl' \
#   && apt-get update -qq \
#   && apt-get install -yq $buildDeps libc6-i386 lib32z1 --no-install-recommends \
#   && echo "===> Install AVG..." \
#   && curl -Ls http://download.avgfree.com/filedir/inst/avg2013flx-r3118-a6926.i386.deb > /tmp/avg.deb \
#   && dpkg -i /tmp/avg.deb \
#   && /etc/init.d/avgd restart \
#   && avgcfgctl -w UpdateVir.sched.Task.Disabled=true \
#   && avgcfgctl -w Default.setup.daemonize=false \
#   && avgcfgctl -w Default.setup.features.antispam=false \
#   && avgcfgctl -w Default.setup.features.oad=false \
#   && avgcfgctl -w Default.setup.features.scheduler=false \
#   && avgcfgctl -w Default.setup.features.tcpd=false \
#   && sed -i 's/Severity=INFO/Severity=None/g' /opt/avg/av/cfg/scand.ini \
#   && sed -i 's/Severity=INFO/Severity=None/g' /opt/avg/av/cfg/tcpd.ini \
#   && sed -i 's/Severity=INFO/Severity=None/g' /opt/avg/av/cfg/wd.ini 


# ARG ZONE_KEY
# ENV ZONE_KEY=$ZONE_KEY

# ENV ZONE 1.3.0

# RUN buildDeps='ca-certificates wget build-essential' \
#   && apt-get update -qq \
#   && apt-get install -yq $buildDeps libc6-i386 \
#   && echo "===> Install Zoner AV..." \
#   # && wget -q -P /tmp http://update.zonerantivirus.com/download/zav-${ZONE}-ubuntu-amd64.deb \
#   && wget --progress=bar:force -P /tmp https://github.com/maliceio/malice-av/raw/master/zoner/zav-1.3.0-debian-amd64.deb \
#   && dpkg -i /tmp/zav-${ZONE}-debian-amd64.deb; \
#   if [ "x$ZONE_KEY" != "x" ]; then \
#   echo "===> Updating License Key..."; \
#   sed -i "s/UPDATE_KEY.*/UPDATE_KEY = ${ZONE_KEY}/g" /etc/zav/zavd.conf; \
#   fi 


# RUN mkdir -p /opt/malice
# RUN if [ "x$ZONE_KEY" != "x" ]; then \
#   echo "===> Update zoner definitions..."; \
#   /etc/init.d/zavd update; \
#   fi


RUN buildDeps='ca-certificates wget' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 --no-install-recommends \
  && set -x \
  && echo "===> Install F-PROT..." \
# COPY fp-Linux.x86.32-ws.tar.gz /tmp/fp-Linux.x86.32-ws.tar.gz
   && wget https://github.com/maliceio/malice-av/raw/master/fprot/fp-Linux.x86.32-ws.tar.gz \
   -O /tmp/fp-Linux.x86.32-ws.tar.gz \
RUN tar -C /opt -zxvf /tmp/fp-Linux.x86.32-ws.tar.gz \
  && ln -fs /opt/f-prot/fpscan /usr/local/bin/fpscan \
  && ln -fs /opt/f-prot/fpscand /usr/local/sbin/fpscand \
  && ln -fs /opt/f-prot/fpmon /usr/local/sbin/fpmon \
  && cp /opt/f-prot/f-prot.conf.default /opt/f-prot/f-prot.conf \
  && ln -fs /opt/f-prot/f-prot.conf /etc/f-prot.conf \
  && chmod a+x /opt/f-prot/fpscan \
  && chmod u+x /opt/f-prot/fpupdate \
  && ln -fs /opt/f-prot/man_pages/scan-mail.pl.8 /usr/share/man/man8/ 





#install clamAv
RUN apt-get install -yq clamav
# RUN apk --update add --no-cache -t .build-deps \
#   build-base \
#   mercurial \
#   musl-dev \
#   openssl \
#   bash \
#   wget \
#   git 

RUN buildDeps='libreadline-dev:i386 \
  ca-certificates \
  libc6-dev:i386 \
  build-essential \
  gcc-multilib \
  cabextract \
  mercurial \
  git-core \
  unzip \
  curl' \
  && set -x \
  && dpkg --add-architecture i386 && apt-get update -qq \
  && apt-get install -y $buildDeps libc6-i386 --no-install-recommends \
  && echo "===> Install taviso/loadlibrary..." \
  && git clone https://github.com/taviso/loadlibrary.git /loadlibrary \
  && echo "===> Download 32-bit antimalware update file.." \
  && curl -L --output /loadlibrary/engine/mpam-fe.exe "https://www.microsoft.com/security/encyclopedia/adlpackages.aspx?arch=x86" \
  && cd /loadlibrary/engine \
  && cabextract mpam-fe.exe \
  && rm mpam-fe.exe \
  && cd /loadlibrary \
  && make -j2 





# ENV LANG en_US.UTF-8
# ENV LANGUAGE en_US:en
# ENV LC_ALL en_US.UTF-8
# ENV TERM=screen-256color
# RUN apt-get update \
#   && apt-get install -yq locales \
#   && locale-gen en_US.UTF-8 

# ARG KASPERSKY_KEY
# ENV KASPERSKY_KEY=$KASPERSKY_KEY
# RUN if [ "x$KASPERSKY_KEY" != "x" ]; then \
#   echo "===> Adding Kaspersky License Key..."; \
#   mkdir -p /etc/kaspersky; \
#   echo -n "$KASPERSKY_KEY" | base64 -d > /etc/kaspersky/license.key ; \
#   fi



# COPY config/docker.conf /etc/kaspersky/docker.conf
# RUN buildDeps='ca-certificates libc6-dev:i386 unzip wget' \
#   && set -x \
#   && dpkg --add-architecture i386 \
#   && apt-get update \
#   && apt-get install -yq $buildDeps libc6-i386 libcurl4-openssl-dev curlftpfs \
#   && echo "===> Install Kaspersky..." \
#   && wget --progress=bar:force https://products.s.kaspersky-labs.com/multilanguage/file_servers/kavlinuxserver8.0/kav4fs_8.0.4-312_i386.deb -P /tmp \
#   && DEBIAN_FRONTEND=noninteractive dpkg --force-architecture -i /tmp/kav4fs_8.0.4-312_i386.deb \
#   && chmod a+s /opt/kaspersky/kav4fs/bin/kav4fs-setup.pl \
#   && chmod a+s /opt/kaspersky/kav4fs/bin/kav4fs-control \
#   && chmod 0777 /etc/kaspersky/license.key \
#   && /opt/kaspersky/kav4fs/bin/kav4fs-control -L --validate-on-install /etc/kaspersky/license.key; sleep 3  \
#   && /opt/kaspersky/kav4fs/bin/kav4fs-control -L --install-on-install /etc/kaspersky/license.key; sleep 3  \
#   && echo "===> Setup Kaspersky..." \
#   && /opt/kaspersky/kav4fs/bin/kav4fs-setup.pl --auto-install=/etc/kaspersky/docker.conf; sleep 10 
#   # && echo "===> Clean up unnecessary files..." \
#   # && apt-get purge -y --auto-remove $buildDeps \
#   # && apt-get clean \
#   # && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/*





ENV ESCAN 7.0-20

RUN buildDeps='wget ca-certificates gdebi' \
  && set -x \
  && dpkg --add-architecture i386 \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 --no-install-recommends \
  && echo "===> Install eScan AV..." \
  && wget -q -P /tmp https://www.microworldsystems.com/download/linux/soho/deb/escan-antivirus-wks.amd64.deb \
  && DEBIAN_FRONTEND=noninteractive gdebi -n /tmp/escan-antivirus-wks.amd64.deb 
  # && echo "===> Clean up unnecessary files..." \
  # && apt-get remove -y $buildDeps \
  # && apt-get clean \
  # && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg



ENV DRWEB 11.0.6

# Install Dr.WEB AV
# COPY drweb-11.0.5-av-linux-amd64.run /tmp/drweb-11.0.5-av-linux-amd64.run
RUN buildDeps='libreadline-dev:i386 \
    ca-certificates \
    libc6-dev:i386 \
    build-essential \
    gcc-multilib \
    cabextract \
    mercurial \
    git-core \
    unzip \
    wget' \
    && set -x \
    && dpkg --add-architecture i386 && apt-get update -qq \
    && apt-get install -y $buildDeps psmisc gnupg libc6-i386 libfontconfig1 libxrender1 libglib2.0-0 libxi6 xauth \
    # && apt-get install -yq libc6-i386 $buildDeps --no-install-recommends \
    && set -x \
    && echo "Install Dr Web..." \
    && cd /tmp \
    && wget --progress=bar:force https://download.geo.drweb.com/pub/drweb/unix/workstation/11.0/drweb-${DRWEB}-av-linux-amd64.run \
    && chmod 755 /tmp/drweb-${DRWEB}-av-linux-amd64.run \
    && DRWEB_NON_INTERACTIVE=yes /tmp/drweb-${DRWEB}-av-linux-amd64.run 
    # && echo "===> Clean up unnecessary files..." \
    # && apt-get purge -y --auto-remove $buildDeps \
    # && apt-get clean \
    # && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/*



# # Install McAfee AV
# RUN set -x \
#     && apt-get update \
#     && apt-get install -yq ca-certificates curl --no-install-recommends \
#     && echo "===> Install McAfee..." \
#     && mkdir -p /usr/local/uvscan  \
#     && curl http://b2b-download.mcafee.com/products/evaluation/vcl/l64/vscl-l64-604-e.tar.gz \
#     |tar -xzf - -C /usr/local/uvscan 
# COPY vscl-l64-604-e/ /usr/local/uvscan 
   



RUN buildDeps='ca-certificates file unzip curl' \
  && dpkg --add-architecture i386 \
  && apt-get update \
  && apt-get install -yq $buildDeps libc6-i386 \
  && echo "===> Install Avira..." \
  && curl -sSL "http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz" \
  | tar -xzf - -C /tmp \
  && mv /tmp/scancl* /opt/avira \
  && curl -sSL -o /tmp/fusebundlegen.zip "http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip" \
  && cd /tmp && unzip /tmp/fusebundlegen.zip \
  && /tmp/fusebundle.bin \
  && mv install/fusebundle-linux_glibc22-int.zip /opt/avira \
  && cd /opt/avira && unzip fusebundle-linux_glibc22-int.zip 



ARG AVIRA_KEY
ENV AVIRA_KEY='FJKDUR-FDJKDI-DFJKDIE-FJKDIE'

# COPY hbedv.key /opt/avira
RUN if [ "x$AVIRA_KEY" != "x" ]; then \
  echo "===> Adding Avira License Key..."; \
  mkdir -p /opt/avira; \
  echo -n "$AVIRA_KEY" | base64 -d > /opt/avira/hbedv.key ; \
  fi


RUN buildDeps='ca-certificates \
  build-essential \
  gdebi-core \
  libssl-dev \
  mercurial \
  git-core \
  wget' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps \
  && echo "===> Install Comodo..." \
  && cd /tmp \
  && wget https://cdn.download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb \
  # && wget http://download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb \
  && DEBIAN_FRONTEND=noninteractive gdebi -n cav-linux_x64.deb \
  && DEBIAN_FRONTEND=noninteractive /opt/COMODO/post_setup.sh \
  && /opt/COMODO/cavdiagnostic 




ADD http://download.comodo.com/av/updates58/sigs/bases/bases.cav /opt/COMODO/scanners/bases.cav




# ARG BDKEY
# ENV BDVERSION 7.7-1

# ENV BDURLPART BitDefender_Antivirus_Scanner_for_Unices/Unix/Current/EN_FR_BR_RO/Linux/
# ENV BDURL https://download.bitdefender.com/SMB/Workstation_Security_and_Management/${BDURLPART}

# RUN buildDeps='ca-certificates wget build-essential' \
#   && apt-get update -qq \
#   && apt-get install -yq $buildDeps psmisc \
#   && set -x \
#   && echo "===> Install Bitdefender..." \
#   && cd /tmp \
#   && wget -q http://download.bitdefender.com/SMB/Workstation_Security_and_Management/BitDefender_Antivirus_Scanner_for_Unices/Unix/Current/EN/FreeBSD/Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run  \
#   && chmod +x /tmp/Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run \
#   && sh /tmp/Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run --check \
#   && echo "===> Making installer noninteractive..." \
#   && sed -i 's/^more LICENSE$/cat  LICENSE/' Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run \
#   && sed -i 's/^CRCsum=.*$/CRCsum="0000000000"/' Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run \
#   && sed -i 's/^MD5=.*$/MD5="00000000000000000000000000000000"/' Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run \
#   && (echo 'accept'; echo 'n') | sh /tmp/Bitdefender-Antivirus-Scanner-for-Unices-8.0.0.freebsd.amd64.run; \
#   if [ "x$BDKEY" != "x" ]; then \
#   echo "===> Updating License..."; \
#   oldkey='^Key =.*$'; \
#   newkey="Key = ${BDKEY}"; \
#   sed -i "s|$oldkey|$newkey|g" /opt/BitDefender-scanner/etc/bdscan.conf; \
#   cat /opt/BitDefender-scanner/etc/bdscan.conf; \
#   fi 



#  RUN mkdir -p /opt/malice && echo "accept" | bdscan --update



COPY aviraupdate.sh /opt/malice/aviraupdate

# Update McAfee Definitions
# COPY mcafeeupdate.sh /usr/local/uvscan/mcafeeupdate

RUN mkdir -p /opt/malice && escan --update


# RUN echo "===> Clean up unnecessary files..." \
#   && apt-get purge -y --auto-remove $buildDeps && apt-get clean \
#   && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg

# Ensure ca-certificates is installed for elasticsearch to use https
#RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates 
  # && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Update AVG Definitions
#  RUN mkdir -p /opt/malice && /etc/init.d/avgd restart && avgupdate

# Update F-PROT Definitions
 RUN mkdir -p /opt/malice && /opt/f-prot/fpupdate






# # Add EICAR Test Virus File to malware folder
# ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

COPY --from=go_builder /bin/avscan /bin/avscan

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]

CMD ["./avscan"]