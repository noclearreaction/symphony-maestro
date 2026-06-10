FROM node:20-slim

# Pin opencode version for reproducibility
RUN npm install -g opencode-ai@1.16.2

# Copy the minimal project fixture into the image
COPY . /app/fixture/

WORKDIR /app/fixture

CMD ["/bin/bash"]
