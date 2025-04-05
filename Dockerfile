# المرحلة الأولى (Stage 1): بناء الصورة
FROM node:16 AS build-stage

WORKDIR /app
COPY . .
RUN npm install

# المرحلة الثانية (Stage 2): تشغيل الصورة
FROM node:16

WORKDIR /app
COPY --from=build-stage /app .
CMD ["npm", "start"]

EXPOSE 80
EXPOSE 443
