FROM node:22-alpine


WORKDIR /usr/src/app

COPY package*.json ./
COPY pnpm-*.yaml ./

RUN corepack enable pnpm
RUN pnpm install --prod

COPY ./dist ./dist
COPY ./proto ./proto
COPY .env ./

RUN chown -R node:node /usr/src/app

USER node

EXPOSE 3000

CMD ["node", "dist/main"]