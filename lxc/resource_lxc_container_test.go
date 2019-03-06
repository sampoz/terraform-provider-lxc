package lxc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"gopkg.in/lxc/go-lxc.v2"
)

func TestLXCContainer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLXCContainerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccLXCContainer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLXCContainerExists(
						t, "lxc_container.accept_test", new(lxc.Container)),
					resource.TestCheckResourceAttr(
						"lxc_container.accept_test", "name", "accept_test"),
				),
			},
		},
	})
}

func testAccCheckLXCContainerExists(t *testing.T, n string, container *lxc.Container) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %v", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config, err := testProviderConfig()
		if err != nil {
			return err
		}

		c := lxc.ActiveContainers(config.LXCPath)
		for i := range c {
			if c[i].Name() == rs.Primary.ID {
				// todo: assigning the container here currently does nothing,
				// we should test some more variables or don't assign it at all
				container = c[i]
				return nil
			}
		}

		return fmt.Errorf("Unable to find running container.")
	}
}

func testAccCheckLXCContainerDestroy(s *terraform.State) error {
	config, err := testProviderConfig()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lxc_container" {
			continue
		}

		c := lxc.ActiveContainers(config.LXCPath)
		for i := range c {
			if c[i].Name() == rs.Primary.ID {
				return fmt.Errorf("Container still exists.")
			}
		}
	}

	return nil
}

var testAccLXCContainer = `
	resource "lxc_container" "accept_test" {
		name = "accept_test"

        network_interface {
            type = "veth"
            options {
               veth.pair = "accept_test_if"
            }
        }
	}`
