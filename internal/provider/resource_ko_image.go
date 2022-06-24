package provider

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/publish"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	baseImage  = "gcr.io/distroless/static:nonroot"
	targetRepo = "gcr.io/jason-chainguard"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceKoBuildCreate,
		ReadContext:   resourceKoBuildRead,
		DeleteContext: resourceKoBuildDelete,

		Schema: map[string]*schema.Schema{
			"importpath": {
				Description: "import path blah",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: func(data interface{}, path cty.Path) diag.Diagnostics {
					// TODO: validate stuff here.
					return nil
				},
				ForceNew: true, // Any time this changes, don't try to update in-place, just create it.
			},
			"image_ref": {
				Description: "image at digest",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func doBuild(ctx context.Context, ip string) (string, error) {
	koDockerRepo := os.Getenv("KO_DOCKER_REPO")

	b, err := build.NewGo(ctx, ".",
		build.WithPlatforms("linux/amd64"), // TODO: needs platforms.
		build.WithBaseImages(func(ctx context.Context, _ string) (name.Reference, build.Result, error) {
			ref := name.MustParseReference(baseImage)
			base, err := remote.Index(ref, remote.WithContext(ctx))
			return ref, base, err
		}),
	)
	if err != nil {
		return "", fmt.Errorf("NewGo: %v", err)
	}
	r, err := b.Build(ctx, ip)
	if err != nil {
		return "", fmt.Errorf("Build: %v", err)
	}

	p, err := publish.NewDefault(koDockerRepo,
		publish.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", fmt.Errorf("NewDefault: %v", err)
	}
	ref, err := p.Publish(ctx, r, ip)
	if err != nil {
		return "", fmt.Errorf("Publish: %v", err)
	}
	return ref.String(), nil
}

const koDockerRepo = "gcr.io/jason-chainguard"

func resourceKoBuildCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ref, err := doBuild(ctx, d.Get("importpath").(string))
	if err != nil {
		return diag.Errorf("doBuild: %v", err)
	}

	d.Set("image_ref", ref)
	d.SetId(ref)
	return nil
}

func resourceKoBuildRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Build the image again, and only unset ID if it changed.

	ref, err := doBuild(ctx, d.Get("importpath").(string))
	if err != nil {
		return diag.Errorf("doBuild: %v", err)
	}

	if ref != d.Id() {
		d.SetId("")
	} else {
		log.Println("image not changed")
	}
	return nil
}

func resourceKoBuildDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO: If we ever want to delete the image from the registry, we can do it here.
	return nil
}