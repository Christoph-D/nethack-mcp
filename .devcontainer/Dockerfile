FROM alpine:3.23

RUN apk add --no-cache bash curl gcompat git go less ripgrep tzdata make nethack tmux \
  gcc ncurses-dev bison flex gzip wget coreutils linux-headers
 
RUN ln -sf /usr/share/zoneinfo/hostlocaltime /etc/localtime

ARG USERNAME=user
RUN addgroup --gid 1000 $USERNAME && adduser --uid 1000 --ingroup $USERNAME --disabled-password $USERNAME
RUN mkdir -p /home/$USER && chown $USERNAME:$USERNAME /home/$USERNAME

USER $USERNAME
WORKDIR /home/$USERNAME

RUN mkdir -p /home/$USERNAME/persistent-config
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/home/$USERNAME/persistent-config/.bash_history" \
  && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"

ENV PATH="$PATH:/usr/local/go/bin:/home/$USERNAME/go/bin/"
RUN go install github.com/Christoph-D/pebbles/cmd/peb@latest

ENV OPENCODE_CONFIG=/home/$USERNAME/persistent-config/opencode.json
RUN curl -fsSL https://opencode.ai/install | bash
ENV PATH="/home/$USERNAME/.opencode/bin:$PATH"

WORKDIR /workspaces/nethack-mcp
