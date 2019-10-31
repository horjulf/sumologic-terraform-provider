package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCollectorCreate,
		Read:   resourceSumologicCollectorRead,
		Delete: resourceSumologicCollectorDelete,
		Update: resourceSumologicCollectorUpdate,
		Exists: resourceSumologicCollectorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Etc/UTC",
			},
			"fields": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"lookup_by_name": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
		},
	}
}

func resourceSumologicCollectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var collector *Collector
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		collector, _ = c.GetCollectorName(d.Id())
		d.SetId(strconv.Itoa(collector.ID))
	} else {
		collector, _ = c.GetCollector(id)
	}

	if collector == nil {
		log.Printf("[WARN] Collector not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)
	d.Set("timezone", collector.TimeZone)
	d.Set("fields", collector.Fields)

	return nil
}

func resourceSumologicCollectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		return c.DeleteCollector(id)
	}

	return nil
}

func resourceSumologicCollectorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		collector, err := c.GetCollectorName(d.Get("name").(string))
		if err != nil {
			return err
		}

		if collector != nil {
			d.SetId(strconv.Itoa(collector.ID))
		}
	}

	if d.Id() == "" {
		id, err := c.CreateCollector(Collector{
			CollectorType: "Hosted",
			Name:          d.Get("name").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCollectorUpdate(d, meta)
}

func resourceSumologicCollectorUpdate(d *schema.ResourceData, meta interface{}) error {

	collector := resourceToCollector(d)

	c := meta.(*Client)
	err := c.UpdateCollector(collector)

	if err != nil {
		return err
	}

	return resourceSumologicCollectorRead(d, meta)
}

func resourceSumologicCollectorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, nil
	}

	collector, _ := c.GetCollector(id)
	if collector != nil {
		return true, nil
	}

	return false, nil
}

func resourceToCollector(d *schema.ResourceData) Collector {
	id, _ := strconv.Atoi(d.Id())

	return Collector{
		ID:            id,
		CollectorType: "Hosted",
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		TimeZone:      d.Get("timezone").(string),
		Fields:        d.Get("fields").(map[string]interface{}),
	}
}
