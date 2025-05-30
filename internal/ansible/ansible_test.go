package ansible

import (
	"context"
	"fmt"
	"maps"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/resource"

	"sigs.k8s.io/yaml"
)

var _ = Describe("Ansible", func() {
	var (
		ctx     context.Context
		listing *resource.ProductListingDeclaration
	)

	BeforeEach(func() {
		ctx = context.TODO()
		listing = &resource.ProductListingDeclaration{}
	})

	When("generating inventory for all components", func() {
		var mapping MappingDeclaration

		BeforeEach(func() {
			mapping = MappingDeclaration{
				ContainerComponents: ComponentCertificationConfig[ContainerCertTarget]{},
				HelmChartComponents: ComponentCertificationConfig[HelmCertTarget]{},
				OperatorComponents:  ComponentCertificationConfig[OperatorCertTarget]{},
			}
		})

		When("the product listing has no componenents defined", func() {
			It("should return an error", func() {
				_, err := GenerateInventory(ctx, listing, mapping)
				Expect(err).To(MatchError(ErrNoComponentsDeclared))
			})
		})

		When("valid components of each type are defined in the listing", func() {
			var container, operator, helmchart resource.Component
			BeforeEach(func() {
				container = resource.Component{
					ID:   "container",
					Type: resource.ComponentTypeContainer,
					Container: &resource.ContainerComponent{
						Type:          resource.ComponentTypeContainer,
						OSContentType: resource.ContentTypeUBI,
					},
				}
				operator = resource.Component{
					ID:   "operator",
					Type: resource.ComponentTypeContainer,
					Container: &resource.ContainerComponent{
						Type:               resource.ComponentTypeContainer,
						OSContentType:      resource.ContentTypeOperatorBundle,
						DistributionMethod: resource.ContainerDistributionExternal,
					},
				}

				helmchart = resource.Component{
					ID:   "helmchart",
					Type: resource.ComponentTypeHelmChart,
					HelmChart: &resource.HelmChartComponent{
						ChartName: "placeholder-chart",
					},
				}

				listing.With.Components = []*resource.Component{
					&container,
					&operator,
					&helmchart,
				}
			})
			When("valid component mappings exist for each component", func() {
				BeforeEach(func() {
					mapping.ContainerComponents[container.ID] = ContainerCertTarget{
						ImageRef: "example.com/example/image",
						Tags:     []ContainerTagAndOptions{{Tag: "tag"}},
					}
					mapping.OperatorComponents[operator.ID] = OperatorCertTarget{
						ImageRef: "example.com/example/image",
						Tags:     []OperatorTagandOptions{{Tag: "tag"}},
					}
					mapping.HelmChartComponents[helmchart.ID] = HelmCertTarget{
						ChartURI: "https://example.com/path/to/chart.0.0.1.tgz",
					}
				})
			})
			It("should be valid YAML", func() {
				inventory, err := GenerateInventory(ctx, listing, mapping)
				Expect(err).ToNot(HaveOccurred())
				Expect(inventory[InventoryKeyContainer]).To(HaveLen(1))
				Expect(inventory[InventoryKeyOperator]).To(HaveLen(1))
				Expect(inventory[InventoryKeyHelmChart]).To(HaveLen(1))

				b, err := yaml.Marshal(inventory)
				Expect(err).ToNot(HaveOccurred())
				Expect(b).ToNot(BeEmpty())
			})
		})
	})

	When("generating inventory for container components", func() {
		var certTarget ComponentCertificationConfig[ContainerCertTarget]

		BeforeEach(func() {
			certTarget = ComponentCertificationConfig[ContainerCertTarget]{}
		})

		When("the product contains no container components", func() {
			It("should produce an empty container component host map", func() {
				hostmap, err := generateContainerComponentInventory(ctx, listing, certTarget)
				Expect(err).ToNot(HaveOccurred())
				Expect(hostmap).To(HaveLen(0))
			})
		})

		When("the product contains a valid application container component", func() {
			var validApplicationContainerComponent resource.Component
			BeforeEach(func() {
				validApplicationContainerComponent = resource.Component{
					ID:   "placeholder",
					Type: resource.ComponentTypeContainer,
					Container: &resource.ContainerComponent{
						Type:          resource.ComponentTypeContainer,
						OSContentType: resource.ContentTypeUBI,
					},
				}

				listing.With.Components = []*resource.Component{&validApplicationContainerComponent}
			})
			When("the mapping does not contain a corresponding entry", func() {
				It("should skip the component and produce an empty host map", func() {
					hostmap, err := generateContainerComponentInventory(ctx, listing, certTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(hostmap).To(HaveLen(0))
				})
			})

			When("the mapping contains a corresponding entry", func() {
				var validContainerCertTarget ContainerCertTarget
				BeforeEach(func() {
					validContainerCertTarget = ContainerCertTarget{
						ImageRef: "example.com/example/image",
						Tags:     []ContainerTagAndOptions{{Tag: "tag"}},
					}

					certTarget[validApplicationContainerComponent.ID] = validContainerCertTarget
				})

				It("should add the container to the hostmap", func() {
					hostmap, err := generateContainerComponentInventory(ctx, listing, certTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(hostmap).To(HaveLen(1))
					imageTarget := strings.Join([]string{validContainerCertTarget.ImageRef, validContainerCertTarget.Tags[0].Tag}, ":")
					hostname := normalizeContainerHostname(
						strings.Join([]string{
							validApplicationContainerComponent.ID,
							imageTarget,
						}, "-"),
					)
					Expect(hostmap).To(HaveKey(hostname))
					Expect(hostmap[hostname].Image).To(Equal(imageTarget))
					Expect(hostmap[hostname].Component.ID).To(Equal(validApplicationContainerComponent.ID))
				})

				When("top level tool flags are defined", func() {
					var topLevelToolFlags map[string]any
					BeforeEach(func() {
						validContainerCertTarget.ToolFlags = map[string]any{}
						topLevelToolFlags = map[string]any{"foo": 3}
						maps.Copy(validContainerCertTarget.ToolFlags, topLevelToolFlags)
						// Re-add the target to the mapping
						certTarget[validApplicationContainerComponent.ID] = validContainerCertTarget
					})

					It("should pass the tool flags to the returned host map", func() {
						hostmap, err := generateContainerComponentInventory(ctx, listing, certTarget)
						Expect(err).ToNot(HaveOccurred())
						Expect(hostmap).To(HaveLen(1))
						imageTarget := strings.Join([]string{validContainerCertTarget.ImageRef, validContainerCertTarget.Tags[0].Tag}, ":")
						hostname := normalizeContainerHostname(
							strings.Join([]string{
								validApplicationContainerComponent.ID,
								imageTarget,
							}, "-"),
						)
						Expect(hostmap[hostname].ToolFlags).To(Equal(topLevelToolFlags))
					})

					When("tool flags are defined on a specific tag", func() {
						var tagSpecificToolFlags map[string]any
						BeforeEach(func() {
							validContainerCertTarget.Tags[0].ToolFlags = map[string]any{}
							tagSpecificToolFlags = map[string]any{"foo": 4}
							maps.Copy(validContainerCertTarget.Tags[0].ToolFlags, tagSpecificToolFlags)
							// Re-add the target to the mapping
							certTarget[validApplicationContainerComponent.ID] = validContainerCertTarget
						})
						It("should override top level flags for the specified tag", func() {
							hostmap, err := generateContainerComponentInventory(ctx, listing, certTarget)
							Expect(err).ToNot(HaveOccurred())
							Expect(hostmap).To(HaveLen(1))
							imageTarget := strings.Join([]string{validContainerCertTarget.ImageRef, validContainerCertTarget.Tags[0].Tag}, ":")
							hostname := normalizeContainerHostname(
								strings.Join([]string{
									validApplicationContainerComponent.ID,
									imageTarget,
								}, "-"),
							)
							Expect(hostmap[hostname].ToolFlags).To(Equal(tagSpecificToolFlags))
						})
					})
				})
			})
		})
	})

	When("generating inventory for operator components", func() {
		var certTarget ComponentCertificationConfig[OperatorCertTarget]

		BeforeEach(func() {
			certTarget = ComponentCertificationConfig[OperatorCertTarget]{}
		})

		When("the product contains no operator components", func() {
			It("should produce an empty host map", func() {
				hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
				Expect(err).ToNot(HaveOccurred())
				Expect(hostmap).To(HaveLen(0))
			})
		})

		When("the product contains a valid operator component", func() {
			var validOperatorComponent resource.Component
			BeforeEach(func() {
				validOperatorComponent = resource.Component{
					ID:   "placeholder",
					Type: resource.ComponentTypeContainer,
					Container: &resource.ContainerComponent{
						Type:               resource.ComponentTypeContainer,
						OSContentType:      resource.ContentTypeOperatorBundle,
						DistributionMethod: resource.ContainerDistributionExternal,
					},
				}

				listing.With.Components = []*resource.Component{&validOperatorComponent}
			})

			When("the mapping does not contain a corresponding entry", func() {
				It("should skip the component and produce an empty host map", func() {
					hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(hostmap).To(HaveLen(0))
				})
			})

			When("the mapping contains a corresponding entry", func() {
				var (
					validCertTarget OperatorCertTarget
					imageTarget     string
					hostname        string
				)
				BeforeEach(func() {
					validCertTarget = OperatorCertTarget{
						ImageRef: "example.com/example/image",
						Tags:     []OperatorTagandOptions{{Tag: "tag"}},
					}

					certTarget[validOperatorComponent.ID] = validCertTarget
				})

				JustBeforeEach(func() {
					imageTarget = strings.Join([]string{validCertTarget.ImageRef, validCertTarget.Tags[0].Tag}, ":")
					hostname = normalizeContainerHostname(
						strings.Join([]string{
							validOperatorComponent.ID,
							imageTarget,
						}, "-"),
					)
				})

				It("should add the operator to the hostmap", func() {
					hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(hostmap).To(HaveLen(1))
					Expect(hostmap).To(HaveKey(hostname))
					Expect(hostmap[hostname].Image).To(Equal(imageTarget))
					Expect(hostmap[hostname].Component.ID).To(Equal(validOperatorComponent.ID))
				})

				When("a top level index image is defined", func() {
					var topLevelIndexImage string
					BeforeEach(func() {
						topLevelIndexImage = "example.io/example/index-image:tag"
						validCertTarget.IndexImage = topLevelIndexImage
						certTarget[validOperatorComponent.ID] = validCertTarget
					})

					It("should pass the index image to entry", func() {
						hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
						Expect(err).ToNot(HaveOccurred())
						Expect(hostmap).To(HaveLen(1))
						Expect(hostmap[hostname].IndexImage).To(Equal(topLevelIndexImage))
					})

					When("a tag-specific index image is defined", func() {
						var tagSpecificIndexImage string
						BeforeEach(func() {
							tagSpecificIndexImage = "example.io/example/index-image:tag-specific"
							validCertTarget.IndexImage = tagSpecificIndexImage
							certTarget[validOperatorComponent.ID] = validCertTarget
						})

						It("should pass the index image to entry", func() {
							hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
							Expect(err).ToNot(HaveOccurred())
							Expect(hostmap).To(HaveLen(1))
							Expect(hostmap[hostname].IndexImage).To(Equal(tagSpecificIndexImage))
						})
					})
				})

				When("top level tool flags are defined", func() {
					var topLevelToolFlags map[string]any
					BeforeEach(func() {
						validCertTarget.ToolFlags = map[string]any{}
						topLevelToolFlags = map[string]any{"foo": 3}
						maps.Copy(validCertTarget.ToolFlags, topLevelToolFlags)
						// Re-add the target to the mapping
						certTarget[validOperatorComponent.ID] = validCertTarget
					})

					It("should pass the tool flags to the returned host map", func() {
						hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
						Expect(err).ToNot(HaveOccurred())
						Expect(hostmap).To(HaveLen(1))
						imageTarget := strings.Join([]string{validCertTarget.ImageRef, validCertTarget.Tags[0].Tag}, ":")
						hostname := normalizeContainerHostname(
							strings.Join([]string{
								validOperatorComponent.ID,
								imageTarget,
							}, "-"),
						)
						Expect(hostmap[hostname].ToolFlags).To(Equal(topLevelToolFlags))
					})

					When("tool flags are defined on a specific tag", func() {
						var tagSpecificToolFlags map[string]any
						BeforeEach(func() {
							validCertTarget.Tags[0].ToolFlags = map[string]any{}
							tagSpecificToolFlags = map[string]any{"foo": 4}
							maps.Copy(validCertTarget.Tags[0].ToolFlags, tagSpecificToolFlags)
							// Re-add the target to the mapping
							certTarget[validOperatorComponent.ID] = validCertTarget
						})
						It("should override top level flags for the specified tag", func() {
							hostmap, err := generateOperatorComponentInventory(ctx, listing, certTarget)
							Expect(err).ToNot(HaveOccurred())
							Expect(hostmap).To(HaveLen(1))
							imageTarget := strings.Join([]string{validCertTarget.ImageRef, validCertTarget.Tags[0].Tag}, ":")
							hostname := normalizeContainerHostname(
								strings.Join([]string{
									validOperatorComponent.ID,
									imageTarget,
								}, "-"),
							)
							Expect(hostmap[hostname].ToolFlags).To(Equal(tagSpecificToolFlags))
						})
					})
				})
			})
		})
	})

	When("generating inventory for helm chart components", func() {
		var certTarget ComponentCertificationConfig[HelmCertTarget]

		BeforeEach(func() {
			certTarget = ComponentCertificationConfig[HelmCertTarget]{}
		})

		When("the product contains no helm chart components", func() {
			It("should produce an empty host map", func() {
				hostmap, err := generateHelmComponentInventory(ctx, listing, certTarget)
				Expect(err).ToNot(HaveOccurred())
				Expect(hostmap).To(HaveLen(0))
			})
		})

		When("the product contains a valid helm chart component", func() {
			var validComponent resource.Component
			BeforeEach(func() {
				validComponent = resource.Component{
					ID:   "placeholder",
					Type: resource.ComponentTypeHelmChart,
					HelmChart: &resource.HelmChartComponent{
						ChartName: "placeholder-chart",
					},
				}

				listing.With.Components = []*resource.Component{&validComponent}
			})

			When("the mapping does not contain a corresponding entry", func() {
				It("should skip the component and produce an empty host map", func() {
					hostmap, err := generateHelmComponentInventory(ctx, listing, certTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(hostmap).To(HaveLen(0))
				})
			})

			When("the mapping contains a corresponding entry", func() {
				var (
					validCertTarget HelmCertTarget
					chartURI        string
					hostname        string
				)
				BeforeEach(func() {
					chartURI = "https://example.com/path/to/chart.0.0.1.tgz"
					validCertTarget = HelmCertTarget{
						ChartURI: chartURI,
					}

					certTarget[validComponent.ID] = validCertTarget
				})

				JustBeforeEach(func() {
					hostname = normalizeChartHostname(
						strings.Join([]string{
							validComponent.ID,
							chartURI,
						}, "-"),
					)
					fmt.Println(hostname)
				})

				It("should add the helm chart to the hostmap", func() {
					hostmap, err := generateHelmComponentInventory(ctx, listing, certTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(hostmap).To(HaveLen(1))
					Expect(hostmap).To(HaveKey(hostname))
					Expect(hostmap[hostname].ChartURI).To(Equal(chartURI))
					Expect(hostmap[hostname].Component.ID).To(Equal(validComponent.ID))
				})

				When("top level tool flags are defined", func() {
					var toolFlags map[string]any
					BeforeEach(func() {
						validCertTarget.ToolFlags = map[string]any{}
						toolFlags = map[string]any{"foo": 3}
						maps.Copy(validCertTarget.ToolFlags, toolFlags)
						// Re-add the target to the mapping
						certTarget[validComponent.ID] = validCertTarget
					})

					It("should pass the tool flags to the returned host map", func() {
						hostmap, err := generateHelmComponentInventory(ctx, listing, certTarget)
						Expect(err).ToNot(HaveOccurred())
						Expect(hostmap).To(HaveLen(1))
						Expect(hostmap[hostname].ToolFlags).To(Equal(toolFlags))
					})
				})
			})
		})
	})
})
