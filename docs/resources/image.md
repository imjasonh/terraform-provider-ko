---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ko_image Resource - terraform-provider-ko"
subcategory: ""
description: |-
  Sample resource in the Terraform provider scaffolding.
---

# ko_image (Resource)

Sample resource in the Terraform provider scaffolding.

## Example Usage

```terraform
resource "ko_image" "example" {
  importpath = "github.com/google/ko"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `importpath` (String) import path to build

### Optional

- `base_image` (String) base image to use
- `platforms` (String) platforms to build
- `working_dir` (String) working directory for the build

### Read-Only

- `id` (String) The ID of this resource.
- `image_ref` (String) built image reference by digest

