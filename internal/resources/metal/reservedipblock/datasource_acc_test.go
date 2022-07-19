package reservedipblock

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceMetalReservedIPBlockConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
	name = "tfacc-reserved_ip_block-%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	metro       = "sv"
	type        = "public_ipv4"
	quantity    = 2
}

data "equinix_metal_reserved_ip_block" "test" {
	project_id  = equinix_metal_project.foobar.id
	ip_address  = cidrhost(equinix_metal_reserved_ip_block.test.cidr_notation,1)
}

data "equinix_metal_reserved_ip_block" "test_id" {
	id  = equinix_metal_reserved_ip_block.test.id
}
`, name)
}

func TestAccDataSourceMetalReservedIPBlock_basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { tfacc.PreCheck(t) },
		Providers:    tfacc.AccProviders,
		CheckDestroy: testAccMetalReservedIPBlockCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalReservedIPBlockConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_reserved_ip_block.test", "id",
						"data.equinix_metal_reserved_ip_block.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_reserved_ip_block.test", "cidr_notation",
						"data.equinix_metal_reserved_ip_block.test_id", "cidr_notation",
					),
				),
			},
		},
	})
}
