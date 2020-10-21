# EKS Support

To create EKS cluster just add or change `controlPlaneProvider` and `bootstrapProvider` to eks in cluster spec.

```yaml
bootstrapProvider:
  name: eks
controlPlaneProvider:
  name: eks
```