package metro

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalMetro_basic(t *testing.T) {
	testMetro := "da"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { tfacc.PreCheck(t) },
		Providers: tfacc.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalMetroConfig_basic(testMetro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_metro.test", "code", testMetro),
				),
			},
			{
				Config: testAccDataSourceMetalMetroConfig_capacityReasonable(testMetro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_metro.test", "code", testMetro),
				),
			},
			{
				Config:      testAccDataSourceMetalMetroConfig_capacityUnreasonable(testMetro),
				ExpectError: matchErrNoCapacity,
			},
			{
				Config:      testAccDataSourceMetalMetroConfig_capacityUnreasonableMultiple(testMetro),
				ExpectError: matchErrNoCapacity,
			},
		},
	})
}

func testAccDataSourceMetalMetroConfig_basic(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
}
`, facCode)
}

func testAccDataSourceMetalMetroConfig_capacityUnreasonable(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1000
    }
}
`, facCode)
}

func testAccDataSourceMetalMetroConfig_capacityReasonable(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1
    }
    capacity {
        plan = "c3.medium.x86"
        quantity = 1
    }
}
`, facCode)
}

func testAccDataSourceMetalMetroConfig_capacityUnreasonableMultiple(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1
    }
    capacity {
        plan = "c3.medium.x86"
        quantity = 1000
    }
}
`, facCode)
}
