# Dockerfile for building your application.
# Defines the final image that contains content from both the image and template.
# Step 1 build the application
FROM golang:1.17.0 as buildEnv
# Get the source code of application
COPY . /project/src/myservice
# Set env variable fo GO environment
ENV GOPATH=/project
# ensure all the dependencies are present
WORKDIR /project/src/myservice
RUN go mod tidy
# build the executable application
#WORKDIR /project
#RUN ls -al /project
RUN go build -o /userapp
# ensure execution priviledges for application
RUN chmod +x /userapp

# Step 2 build the final image
# use 
FROM registry.access.redhat.com/ubi8/ubi-minimal
WORKDIR /
# get the application from the build container (buildEnv)
COPY --from=buildEnv /userapp ./
RUN mkdir -p /var/app
RUN chmod 777 /var/app
# expose the port used
ENV PORT=8080
EXPOSE ${PORT}
#define label
LABEL "version"="1.0.0" \
      "maintainer"="jerome.tarte@fr.ibm.com" \
      "description"="provide a sample application written in GO lang"
# Pass control your application
CMD ["/userapp"]