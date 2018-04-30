# Nginx status monitor
This example visualize `requests` written in stub_status

## Requirements
If your nginx was built from source, let's check config as bellow
```
configure arguments: \
	 ...
	 --with-http_stub_status_module \
	 ...
```

And add next directive in nginx.conf.
```
location /nginx_status {
    stub_status on;
    access_log off;
    //change access policy for your use
    //allow 127.0.0.1;
    //deny all;
}
```

## Run

``` sh
>> go run nginx_status.go --target http://$FIRST_HOST/nginx_status --target http://$SECOND_HOST/nginx_status
```
or
``` sh
>> go build nginx_status.go
>> nginx_status --target http://$FIRST_HOST/nginx_status --target http://$SECOND_HOST/nginx_status
```
