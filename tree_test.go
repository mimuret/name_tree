package name_tree_test

import (
	"testing"

	name_tree "github.com/mimuret/name_tree"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Hoge struct {
	V string
}

func TestNode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dnsutils Suite")
}

var _ = Describe("tree", func() {
	Context("NewTree", func() {
		var (
			v  = &Hoge{"value"}
			tr *name_tree.Tree[Hoge]
		)
		BeforeEach(func() {
			tr = name_tree.NewTree("examplE.Jp", v)
		})
		It("returns node", func() {
			Expect(tr).NotTo(BeNil())
		})
	})
	Context("Tree", func() {
		var (
			nullVar *Hoge
			err     error
			root    *name_tree.Tree[Hoge]
		)
		BeforeEach(func() {
			root = name_tree.NewTree("examplE.Jp", &Hoge{"root"})
		})
		Context("InsertNode", func() {
			When("target is not subdomain", func() {
				BeforeEach(func() {
					err = root.InsertNode(name_tree.NewNode("example.com", nullVar))
				})
				It("returnes error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
			When("target is same domain", func() {
				BeforeEach(func() {
					err = root.InsertNode(name_tree.NewNode("example.jp", nullVar))
				})
				It("returnes error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
			When("target is sub domain", func() {
				BeforeEach(func() {
					err = root.InsertNode(name_tree.NewNode("sub.example.jp", nullVar))
				})
				It("returnes error", func() {
					Expect(err).To(Succeed())
				})
			})
			When("target is sub domain", func() {
				BeforeEach(func() {
					err = root.InsertNode(name_tree.NewNode("sub1.ent.example.jp", nullVar))
				})
				It("returnes error", func() {
					Expect(err).To(Succeed())
				})
			})
		})
		Context("GetNode", func() {
			var (
				sub   *name_tree.Node[Hoge]
				exist bool
			)
			BeforeEach(func() {
				err = root.InsertNode(name_tree.NewNode("sub1.ent.example.jp", &Hoge{"sub1"}))
			})
			When("name not exist", func() {
				BeforeEach(func() {
					sub, exist = root.GetNode("example.net.")
				})
				It("returnes nil", func() {
					Expect(exist).To(BeFalse())
					Expect(sub).To(BeNil())
				})
			})
			When("name exist", func() {
				BeforeEach(func() {
					sub, exist = root.GetNode("ent.example.jp.")
				})
				It("returnes nil", func() {
					Expect(exist).To(BeTrue())
					Expect(sub).NotTo(BeNil())
					Expect(sub.GetName()).To(Equal("ent.example.jp."))
					Expect(sub.GetVault()).To(BeNil())
				})
			})
			When("name is out of tree, but subdomian", func() {
				BeforeEach(func() {
					sub, exist = root.GetNode("sub.sub1.ent.example.jp.")
				})
				It("returnes neer node", func() {
					Expect(exist).To(BeFalse())
					Expect(sub).NotTo(BeNil())
					Expect(sub.GetName()).To(Equal("sub1.ent.example.jp."))
					Expect(sub.GetVault()).To(Equal(&Hoge{"sub1"}))
				})
			})
		})
		Context("RemoveNode", func() {
			var (
				sub   *name_tree.Node[Hoge]
				exist bool
			)
			BeforeEach(func() {
				err = root.InsertNode(name_tree.NewNode("sub1.ent.example.jp", &Hoge{"sub1"}))
				Expect(err).To(Succeed())
				err = root.InsertNode(name_tree.NewNode("sub2.ent.example.jp", &Hoge{"sub2"}))
				Expect(err).To(Succeed())
			})
			When("name is out of tree", func() {
				BeforeEach(func() {
					root.RemoveNode("example.net.")
					sub, exist = root.GetNode("example.jp.")
				})
				It("noting to do", func() {
					Expect(exist).To(BeTrue())
				})
			})
			When("name is root", func() {
				BeforeEach(func() {
					root.RemoveNode("example.jp.")
					sub, exist = root.GetNode("example.jp.")
				})
				It("noting to do", func() {
					Expect(exist).To(BeTrue())
				})
			})
			When("name has children", func() {
				BeforeEach(func() {
					root.RemoveNode("ent.example.jp.")
				})
				It("remove node with children", func() {
					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
					sub, exist = root.GetNode("sub1.ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
					sub, exist = root.GetNode("sub2.ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
				})
			})
			When("name has parent ENT node", func() {
				BeforeEach(func() {
					root.RemoveNode("sub1.ent.example.jp.")
				})
				It("remove node, not remove parent ENT node", func() {
					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("ent.example.jp."))
					sub, exist = root.GetNode("sub1.ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("ent.example.jp."))
					sub, exist = root.GetNode("sub2.ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("sub2.ent.example.jp."))
				})
			})
			When("name is out of tree, but subdomian", func() {
				BeforeEach(func() {
					root.RemoveNode("sub.sub1.ent.example.jp.")
				})
				It("nothing to do", func() {
					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("ent.example.jp."))

					sub, exist = root.GetNode("sub1.ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("sub1.ent.example.jp."))

					sub, exist = root.GetNode("sub2.ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("sub2.ent.example.jp."))
				})
			})
		})
		Context("RemoveNodeWithENT", func() {
			var (
				sub   *name_tree.Node[Hoge]
				exist bool
			)
			BeforeEach(func() {
				err = root.InsertNode(name_tree.NewNode("sub1.ent.example.jp", &Hoge{"sub1"}))
				Expect(err).To(Succeed())
				err = root.InsertNode(name_tree.NewNode("sub2.ent.example.jp", &Hoge{"sub2"}))
				Expect(err).To(Succeed())
				err = root.InsertNode(name_tree.NewNode("sub3.ent3.example.jp", &Hoge{"sub3"}))
				Expect(err).To(Succeed())
			})
			When("name is out of tree", func() {
				BeforeEach(func() {
					root.RemoveNodeWithENT("example.net.", nil)
					sub, exist = root.GetNode("example.jp.")
				})
				It("noting to do", func() {
					Expect(exist).To(BeTrue())
				})
			})
			When("name is root", func() {
				BeforeEach(func() {
					root.RemoveNodeWithENT("example.jp.", nil)
					sub, exist = root.GetNode("example.jp.")
				})
				It("noting to do", func() {
					Expect(exist).To(BeTrue())
				})
			})
			When("name has children", func() {
				BeforeEach(func() {
					root.RemoveNodeWithENT("ent.example.jp.", nil)
				})
				It("remove node with children", func() {
					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
					sub, exist = root.GetNode("sub1.ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
					sub, exist = root.GetNode("sub2.ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
				})
			})
			When("name has ENT parent node", func() {
				BeforeEach(func() {
					root.RemoveNodeWithENT("sub3.ent3.example.jp.", nil)
				})
				It("remove with parent ENT node", func() {
					sub, exist = root.GetNode("ent3.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
					sub, exist = root.GetNode("ent3.ent3.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
				})
			})
			When("names has ENT parent node", func() {
				BeforeEach(func() {
					root.RemoveNodeWithENT("sub1.ent.example.jp.", nil)
				})
				It("remove with parent ENT node", func() {
					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("ent.example.jp."))

					root.RemoveNodeWithENT("sub2.ent.example.jp.", nil)

					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeFalse())
					Expect(sub.GetName()).To(Equal("example.jp."))
				})
			})
			When("name is out of tree, but subdomian", func() {
				BeforeEach(func() {
					root.RemoveNode("sub.sub1.ent.example.jp.")
				})
				It("nothing to do", func() {
					sub, exist = root.GetNode("ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("ent.example.jp."))

					sub, exist = root.GetNode("sub1.ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("sub1.ent.example.jp."))

					sub, exist = root.GetNode("sub2.ent.example.jp.")
					Expect(exist).To(BeTrue())
					Expect(sub.GetName()).To(Equal("sub2.ent.example.jp."))
				})
			})
		})
	})
})
