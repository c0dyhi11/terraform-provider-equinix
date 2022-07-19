package gateway

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalGateway_privateIPv4(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { tfacc.PreCheck(t) },
		Providers:    tfacc.AccProviders,
		CheckDestroy: testAccMetalGatewayCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalGatewayConfig_privateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.equinix_metal_gateway.test", "project_id",
						"equinix_metal_project.test", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_gateway.test", "private_ipv4_subnet_size", "8"),
				),
			},
		},
	})
}

func testAccDataSourceMetalGatewayConfig_privateIPv4() string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-gateway-test"
}

resource "equinix_metal_vlan" "test" {
    description = "tfacc-vlan test VLAN in SV"
    metro       = "sv"
    project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_gateway" "test" {
    project_id               = equinix_metal_project.test.id
    vlan_id                  = equinix_metal_vlan.test.id
    private_ipv4_subnet_size = 8
}

data "equinix_metal_gateway" "test" {
    gateway_id = equinix_metal_gateway.test.id
}
`)
}
