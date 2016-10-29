all: wbinflux

wbinflux: *.go
	wbdev user bash -c 'CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags -s'

image: wbinflux
	docker build -t gcr.io/i4-kube-cluster/wbinflux .

push: image
	gcloud docker push gcr.io/i4-kube-cluster/wbinflux
