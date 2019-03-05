package lxc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"gopkg.in/lxc/go-lxc.v2"
)

func TestLXCClone(t *testing.T) {
	testResource(t, testAccLXCClone)
}

func TestLXCCloneSnapshot(t *testing.T) {
	testResource(t, testAccLXCCloneSnapshot)
}

func TestLXCCloneMixed(t *testing.T) {
	testResource(t, testAccLXCCloneMixed)
}

func testResource(t *testing.T, conf string) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLXCCloneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
				Check:  resource.ComposeTestCheckFunc(testCheckLXCCloneResources(t)...),
			},
		},
	})
}

func testCheckLXCCloneResources(t *testing.T) []resource.TestCheckFunc {
	tests := []resource.TestCheckFunc{}
	for i := 0; i <= 5; i++ {
		r := fmt.Sprintf("accept_clone%d", i)
		rId := fmt.Sprintf("lxc_clone.%s", r)
		tests = append(tests, testAccCheckLXCCloneExists(t, rId, new(lxc.Container)))
		tests = append(tests, resource.TestCheckResourceAttr(rId, "name", r))
	}

	return tests
}

func testAccCheckLXCCloneExists(t *testing.T, n string, container *lxc.Container) resource.TestCheckFunc {
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

func testAccCheckLXCCloneDestroy(s *terraform.State) error {
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

var testAccLXCClone = `
	resource "lxc_container" "accept_test" {
		name = "accept_test"
	}
	resource "lxc_clone" "accept_clone0" {
		name = "accept_clone0"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone1" {
		name = "accept_clone1"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone2" {
		name = "accept_clone2"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone3" {
		name = "accept_clone3"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone4" {
		name = "accept_clone4"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone5" {
		name = "accept_clone5"
		source = "${lxc_container.accept_test.name}"
	}
`

var testAccLXCCloneSnapshot = `
	resource "lxc_container" "accept_test" {
		name = "accept_test"
	}
	resource "lxc_clone" "accept_clone0" {
		name = "accept_clone0"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone1" {
		name = "accept_clone1"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
	resource "lxc_clone" "accept_clone2" {
		name = "accept_clone2"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
	resource "lxc_clone" "accept_clone3" {
		name = "accept_clone3"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
	resource "lxc_clone" "accept_clone4" {
		name = "accept_clone4"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
	resource "lxc_clone" "accept_clone5" {
		name = "accept_clone5"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
`

var testAccLXCCloneMixed = `
	resource "lxc_container" "accept_test" {
		name = "accept_test"
	}
	resource "lxc_clone" "accept_clone0" {
		name = "accept_clone0"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone1" {
		name = "accept_clone1"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone2" {
		name = "accept_clone2"
		source = "${lxc_container.accept_test.name}"
	}
	resource "lxc_clone" "accept_clone3" {
		name = "accept_clone3"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
	resource "lxc_clone" "accept_clone4" {
		name = "accept_clone4"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
	resource "lxc_clone" "accept_clone5" {
		name = "accept_clone5"
		source = "${lxc_container.accept_test.name}"
        backend = "overlayfs"
        snapshot = true
	}
`
