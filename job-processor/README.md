# Prometheus Example - Job Processor

This Job Processor example provides a complete working demonstartion that complements this blog post on
[monitoring Go applications with Prometheus](https://scot.coffee/2018/12/monitoring-go-applications-with-prometheus/).

The releases under the `v1.0.0` major version provide a walk-through using commits of how to add these metrics to your
services. For example, you can run the following command to see the changes required to add a metric for tracking the
total number of jobs processed by the workers.

```shell
$ git diff v1.1.0..v1.2.0
```

To run the full example, check out the `v1.4.0` version of the codebase and start the go service using the following
command.

```shell
$ go run job-processor/main.go
```

Then, you can start the demo Grafana and Prometheus server provided by the example by using
[docker compose](https://docker.com). To start the servers, execute the following commands.

```shell
$ cd job-processor
$ docker-compose up
```

You should then be able to open Grafana in your browser by visiting [http://localhost:3000](http://localhost:3000). Once
you have opened Grafana, you can login using `admin` as the username and `admin` as the password. Then, you can view the
example dashboard by opening the **Job Processor** dashboard.
