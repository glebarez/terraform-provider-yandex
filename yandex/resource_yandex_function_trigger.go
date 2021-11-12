package yandex

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/serverless/triggers/v1"
)

const (
	yandexFunctionTriggerDefaultTimeout = 5 * time.Minute

	triggerTypeIoT           = "iot"
	triggerTypeMessageQueue  = "message_queue"
	triggerTypeObjectStorage = "object_storage"
	triggerTypeTimer         = "timer"
	triggerTypeLogGroup      = "log_group"
	triggerTypeLogging       = "logging"
)

var functionTriggerTypesList = []string{
	triggerTypeIoT,
	triggerTypeMessageQueue,
	triggerTypeObjectStorage,
	triggerTypeTimer,
	triggerTypeLogGroup,
	triggerTypeLogging,
}

var levelNameToEnum = map[string]logging.LogLevel_Level{
	"debug": logging.LogLevel_DEBUG,
	"error": logging.LogLevel_ERROR,
	"fatal": logging.LogLevel_FATAL,
	"info":  logging.LogLevel_INFO,
	"trace": logging.LogLevel_TRACE,
	"warn":  logging.LogLevel_WARN,
}

var levelEnumToName = map[logging.LogLevel_Level]string{
	logging.LogLevel_DEBUG: "debug",
	logging.LogLevel_ERROR: "error",
	logging.LogLevel_FATAL: "fatal",
	logging.LogLevel_INFO:  "info",
	logging.LogLevel_TRACE: "trace",
	logging.LogLevel_WARN:  "warn",
}

func resourceYandexFunctionTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceYandexFunctionTriggerCreate,
		Read:   resourceYandexFunctionTriggerRead,
		Update: resourceYandexFunctionTriggerUpdate,
		Delete: resourceYandexFunctionTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(yandexFunctionTriggerDefaultTimeout),
			Update: schema.DefaultTimeout(yandexFunctionTriggerDefaultTimeout),
			Delete: schema.DefaultTimeout(yandexFunctionTriggerDefaultTimeout),
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"function": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"service_account_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"tag": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"retry_attempts": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"retry_interval": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"folder_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			triggerTypeIoT: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: functionTriggerConflictingTypes(triggerTypeIoT),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"registry_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"device_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"topic": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			triggerTypeMessageQueue: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: functionTriggerConflictingTypes(triggerTypeMessageQueue),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"service_account_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"batch_cutoff": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"batch_size": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"visibility_timeout": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			triggerTypeObjectStorage: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: functionTriggerConflictingTypes(triggerTypeObjectStorage),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"suffix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"create": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},

						"update": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},

						"delete": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			triggerTypeTimer: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: functionTriggerConflictingTypes(triggerTypeTimer),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_expression": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			triggerTypeLogGroup: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: functionTriggerConflictingTypes(triggerTypeLogGroup),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_group_ids": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							MinItems: 1,
						},

						"batch_cutoff": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"batch_size": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			triggerTypeLogging: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: functionTriggerConflictingTypes(triggerTypeLogging),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"resource_ids": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							MinItems: 0,
						},

						"resource_types": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							MinItems: 0,
						},

						"levels": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
							MinItems: 0,
						},

						"batch_cutoff": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"batch_size": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"dlq": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"service_account_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceYandexFunctionTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutCreate))
	defer cancel()

	labels, err := expandLabels(d.Get("labels"))
	if err != nil {
		return fmt.Errorf("Error expanding labels while creating Yandex Cloud Functions Trigger: %s", err)
	}

	folderID, err := getFolderID(d, config)
	if err != nil {
		return fmt.Errorf("Error getting folder ID while creating Yandex Cloud Functions Trigger: %s", err)
	}

	req := triggers.CreateTriggerRequest{
		FolderId:    folderID,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      labels,
	}

	retrySettings, err := expandRetrySettings(d)
	if err != nil {
		return err
	}

	dlqSettings, err := expandDLQSettings(d)
	if err != nil {
		return err
	}

	var getInvokeFunctionWithRetry = func() *triggers.InvokeFunctionWithRetry {
		return &triggers.InvokeFunctionWithRetry{
			FunctionId:       d.Get("function.0.id").(string),
			FunctionTag:      d.Get("function.0.tag").(string),
			ServiceAccountId: d.Get("function.0.service_account_id").(string),
			RetrySettings:    retrySettings,
			DeadLetterQueue:  dlqSettings,
		}
	}

	var getInvokeFunctionOnce = func() *triggers.InvokeFunctionOnce {
		return &triggers.InvokeFunctionOnce{
			FunctionId:       d.Get("function.0.id").(string),
			FunctionTag:      d.Get("function.0.tag").(string),
			ServiceAccountId: d.Get("function.0.service_account_id").(string),
		}
	}

	triggerCnt := 0
	if _, ok := d.GetOk(triggerTypeIoT); ok {
		triggerCnt++
		iot := &triggers.Trigger_Rule_IotMessage{
			IotMessage: &triggers.Trigger_IoTMessage{
				RegistryId: d.Get("iot.0.registry_id").(string),
				DeviceId:   d.Get("iot.0.device_id").(string),
				MqttTopic:  d.Get("iot.0.topic").(string),
				Action: &triggers.Trigger_IoTMessage_InvokeFunction{
					InvokeFunction: getInvokeFunctionWithRetry(),
				},
			},
		}

		req.Rule = &triggers.Trigger_Rule{Rule: iot}
	}

	if _, ok := d.GetOk(triggerTypeMessageQueue); ok {
		triggerCnt++
		batch, err := expandBatchSettings(d, "message_queue.0")
		if err != nil {
			return err
		}

		messageQueue := &triggers.Trigger_MessageQueue{
			QueueId:          d.Get("message_queue.0.queue_id").(string),
			ServiceAccountId: d.Get("message_queue.0.service_account_id").(string),
			BatchSettings:    batch,
			Action: &triggers.Trigger_MessageQueue_InvokeFunction{
				InvokeFunction: getInvokeFunctionOnce(),
			},
		}

		if _, ok := d.GetOk("message_queue.0.visibility_timeout"); ok {
			timeout, err := strconv.ParseInt(d.Get("message_queue.0.visibility_timeout").(string), 10, 64)
			if err != nil {
				return fmt.Errorf("Cannot define message_queue.visibility_timeout for Yandex Cloud Functions Trigger: %s", err)
			}
			messageQueue.VisibilityTimeout = &duration.Duration{Seconds: timeout}
		}

		messageQueueRule := &triggers.Trigger_Rule_MessageQueue{MessageQueue: messageQueue}
		req.Rule = &triggers.Trigger_Rule{Rule: messageQueueRule}
	}

	if _, ok := d.GetOk(triggerTypeObjectStorage); ok {
		triggerCnt++

		events := make([]triggers.Trigger_ObjectStorageEventType, 0)
		eventsName := map[string]triggers.Trigger_ObjectStorageEventType{
			"object_storage.0.create": triggers.Trigger_OBJECT_STORAGE_EVENT_TYPE_CREATE_OBJECT,
			"object_storage.0.update": triggers.Trigger_OBJECT_STORAGE_EVENT_TYPE_UPDATE_OBJECT,
			"object_storage.0.delete": triggers.Trigger_OBJECT_STORAGE_EVENT_TYPE_DELETE_OBJECT,
		}
		for k, v := range eventsName {
			if d.Get(k).(bool) {
				events = append(events, v)
			}
		}

		storageTrigger := &triggers.Trigger_ObjectStorage{
			BucketId:  d.Get("object_storage.0.bucket_id").(string),
			Prefix:    d.Get("object_storage.0.prefix").(string),
			Suffix:    d.Get("object_storage.0.suffix").(string),
			EventType: events,
			Action: &triggers.Trigger_ObjectStorage_InvokeFunction{
				InvokeFunction: getInvokeFunctionWithRetry(),
			},
		}

		storageRule := &triggers.Trigger_Rule_ObjectStorage{ObjectStorage: storageTrigger}
		req.Rule = &triggers.Trigger_Rule{Rule: storageRule}
	}

	if _, ok := d.GetOk(triggerTypeTimer); ok {
		triggerCnt++

		timer := triggers.Trigger_Timer{
			CronExpression: d.Get("timer.0.cron_expression").(string),
		}

		if retrySettings != nil || dlqSettings != nil {
			timer.Action = &triggers.Trigger_Timer_InvokeFunctionWithRetry{
				InvokeFunctionWithRetry: getInvokeFunctionWithRetry(),
			}
		} else {
			timer.Action = &triggers.Trigger_Timer_InvokeFunction{
				InvokeFunction: getInvokeFunctionOnce(),
			}
		}

		timerRule := &triggers.Trigger_Rule_Timer{Timer: &timer}
		req.Rule = &triggers.Trigger_Rule{Rule: timerRule}
	}

	if _, ok := d.GetOk(triggerTypeLogGroup); ok {
		triggerCnt++

		cloudLogs := &triggers.Trigger_CloudLogs{
			LogGroupId: convertStringSet(d.Get("log_group.0.log_group_ids").(*schema.Set)),
			Action: &triggers.Trigger_CloudLogs_InvokeFunction{
				InvokeFunction: getInvokeFunctionWithRetry(),
			},
		}
		batch, err := expandBatchSettings(d, "log_group.0")
		if err != nil {
			return err
		}
		if batch != nil {
			cloudLogs.BatchSettings = &triggers.CloudLogsBatchSettings{
				Size:   batch.Size,
				Cutoff: batch.Cutoff,
			}
		}
		req.Rule = &triggers.Trigger_Rule{
			Rule: &triggers.Trigger_Rule_CloudLogs{CloudLogs: cloudLogs},
		}
	}

	if _, ok := d.GetOk(triggerTypeLogging); ok {
		triggerCnt++

		levels := []logging.LogLevel_Level{}

		for _, l := range convertStringSet(d.Get("logging.0.levels").(*schema.Set)) {
			if v, ok := levelNameToEnum[strings.ToLower(l)]; ok {
				levels = append(levels, v)
			}
		}

		logging := &triggers.Trigger_Logging{
			LogGroupId:   d.Get("logging.0.group_id").(string),
			ResourceId:   convertStringSet(d.Get("logging.0.resource_ids").(*schema.Set)),
			ResourceType: convertStringSet(d.Get("logging.0.resource_types").(*schema.Set)),
			Levels:       levels,

			Action: &triggers.Trigger_Logging_InvokeFunction{
				InvokeFunction: getInvokeFunctionWithRetry(),
			},
		}
		batch, err := expandBatchSettings(d, "logging.0")
		if err != nil {
			return err
		}
		if batch != nil {
			logging.BatchSettings = &triggers.LoggingBatchSettings{
				Size:   batch.Size,
				Cutoff: batch.Cutoff,
			}
		}
		req.Rule = &triggers.Trigger_Rule{
			Rule: &triggers.Trigger_Rule_Logging{Logging: logging},
		}
	}

	if triggerCnt != 1 {
		return fmt.Errorf("Yandex Cloud Functions Trigger must have only one any iot, message_queue, object_storage, timer, logging section")
	}

	op, err := config.sdk.WrapOperation(config.sdk.Serverless().Triggers().Trigger().Create(ctx, &req))
	if err != nil {
		return fmt.Errorf("Error while requesting API to create Yandex Cloud Functions Trigger: %s", err)
	}

	protoMetadata, err := op.Metadata()
	if err != nil {
		return fmt.Errorf("Error while requesting API to create Yandex Cloud Function Trigger: %s", err)
	}

	md, ok := protoMetadata.(*triggers.CreateTriggerMetadata)
	if !ok {
		return fmt.Errorf("Could not get Yandex Cloud Functions Trigger ID from create operation metadata")
	}

	d.SetId(md.TriggerId)

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error while requesting API to create Yandex Cloud Functions Trigger: %s", err)
	}

	return resourceYandexFunctionTriggerRead(d, meta)
}

func resourceYandexFunctionTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutCreate))
	defer cancel()

	labels, err := expandLabels(d.Get("labels"))
	if err != nil {
		return fmt.Errorf("Error expanding labels while creating Yandex Cloud Functions Trigger: %s", err)
	}

	d.Partial(true)

	var updatePaths []string
	if d.HasChange("name") {
		updatePaths = append(updatePaths, "name")
	}

	if d.HasChange("description") {
		updatePaths = append(updatePaths, "description")
	}

	if d.HasChange("labels") {
		updatePaths = append(updatePaths, "labels")
	}

	if len(updatePaths) != 0 {
		req := triggers.UpdateTriggerRequest{
			TriggerId:   d.Id(),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Labels:      labels,
			UpdateMask:  &field_mask.FieldMask{Paths: updatePaths},
		}

		op, err := config.sdk.Serverless().Triggers().Trigger().Update(ctx, &req)
		err = waitOperation(ctx, config, op, err)
		if err != nil {
			return fmt.Errorf("Error while requesting API to update Yandex Cloud Functions Trigger: %s", err)
		}

	}

	d.Partial(false)

	return resourceYandexFunctionTriggerRead(d, meta)
}

func flattenYandexFunctionTriggerInvokeOnce(d *schema.ResourceData, function *triggers.InvokeFunctionOnce) error {
	f := map[string]interface{}{
		"id":                 function.FunctionId,
		"tag":                function.FunctionTag,
		"service_account_id": function.ServiceAccountId,
	}
	return d.Set("function", []map[string]interface{}{f})
}

func flattenYandexFunctionTriggerInvokeWithRetry(d *schema.ResourceData, function *triggers.InvokeFunctionWithRetry) error {
	f := map[string]interface{}{
		"id":                 function.FunctionId,
		"tag":                function.FunctionTag,
		"service_account_id": function.ServiceAccountId,
	}

	if retrySettings := function.GetRetrySettings(); retrySettings != nil {
		f["retry_attempts"] = strconv.FormatInt(retrySettings.RetryAttempts, 10)
		if retrySettings.Interval != nil {
			f["retry_interval"] = strconv.FormatInt(retrySettings.Interval.Seconds, 10)
		}
	}

	err := d.Set("function", []map[string]interface{}{f})
	if err != nil {
		return err
	}

	if deadLetter := function.GetDeadLetterQueue(); deadLetter != nil {
		dlq := map[string]interface{}{
			"queue_id":           deadLetter.QueueId,
			"service_account_id": deadLetter.ServiceAccountId,
		}

		err = d.Set("dlq", []map[string]interface{}{dlq})
		if err != nil {
			return err
		}
	}
	return nil
}

func flattenYandexFunctionTrigger(d *schema.ResourceData, trig *triggers.Trigger) error {
	d.Set("name", trig.Name)
	d.Set("folder_id", trig.FolderId)
	d.Set("description", trig.Description)
	d.Set("created_at", getTimestamp(trig.CreatedAt))
	if err := d.Set("labels", trig.Labels); err != nil {
		return err
	}

	if iot := trig.GetRule().GetIotMessage(); iot != nil {
		i := map[string]interface{}{
			"registry_id": iot.RegistryId,
			"device_id":   iot.DeviceId,
			"topic":       iot.MqttTopic,
		}

		err := d.Set(triggerTypeIoT, []map[string]interface{}{i})
		if err != nil {
			return err
		}

		if function := iot.GetInvokeFunction(); function != nil {
			err = flattenYandexFunctionTriggerInvokeWithRetry(d, function)
			if err != nil {
				return err
			}
		}
	} else if messageQueue := trig.GetRule().GetMessageQueue(); messageQueue != nil {
		m := map[string]interface{}{
			"queue_id":           messageQueue.QueueId,
			"service_account_id": messageQueue.ServiceAccountId,
		}

		if messageQueue.VisibilityTimeout != nil {
			m["visibility_timeout"] = strconv.FormatInt(messageQueue.VisibilityTimeout.Seconds, 10)
		}

		if batch := messageQueue.GetBatchSettings(); batch != nil {
			m["batch_size"] = strconv.FormatInt(batch.Size, 10)
			m["batch_cutoff"] = strconv.FormatInt(batch.Cutoff.Seconds, 10)
		}

		err := d.Set(triggerTypeMessageQueue, []map[string]interface{}{m})
		if err != nil {
			return err
		}

		if function := messageQueue.GetInvokeFunction(); function != nil {
			err = flattenYandexFunctionTriggerInvokeOnce(d, function)
			if err != nil {
				return err
			}
		}
	} else if storage := trig.GetRule().GetObjectStorage(); storage != nil {
		s := map[string]interface{}{
			"bucket_id": storage.BucketId,
			"prefix":    storage.Prefix,
			"suffix":    storage.Suffix,
		}

		events := map[triggers.Trigger_ObjectStorageEventType]string{
			triggers.Trigger_OBJECT_STORAGE_EVENT_TYPE_CREATE_OBJECT: "create",
			triggers.Trigger_OBJECT_STORAGE_EVENT_TYPE_UPDATE_OBJECT: "update",
			triggers.Trigger_OBJECT_STORAGE_EVENT_TYPE_DELETE_OBJECT: "delete",
		}

		for _, t := range storage.EventType {
			if _, ok := events[t]; ok {
				s[events[t]] = true
			}
		}

		err := d.Set(triggerTypeObjectStorage, []map[string]interface{}{s})
		if err != nil {
			return err
		}

		if function := storage.GetInvokeFunction(); function != nil {
			err = flattenYandexFunctionTriggerInvokeWithRetry(d, function)
			if err != nil {
				return err
			}
		}
	} else if timer := trig.GetRule().GetTimer(); timer != nil {
		t := map[string]interface{}{
			"cron_expression": timer.CronExpression,
		}

		err := d.Set(triggerTypeTimer, []map[string]interface{}{t})
		if err != nil {
			return err
		}

		if function := timer.GetInvokeFunctionWithRetry(); function != nil {
			err = flattenYandexFunctionTriggerInvokeWithRetry(d, function)
			if err != nil {
				return err
			}
		} else if function := timer.GetInvokeFunction(); function != nil {
			err = flattenYandexFunctionTriggerInvokeOnce(d, function)
			if err != nil {
				return err
			}
		}
	} else if logGroup := trig.GetRule().GetCloudLogs(); logGroup != nil {

		groupIDs := &schema.Set{F: schema.HashString}
		for _, groupID := range logGroup.LogGroupId {
			groupIDs.Add(groupID)
		}
		lg := map[string]interface{}{
			"log_group_ids": groupIDs,
		}
		if batch := logGroup.GetBatchSettings(); batch != nil {
			lg["batch_size"] = strconv.FormatInt(batch.Size, 10)
			lg["batch_cutoff"] = strconv.FormatInt(batch.Cutoff.Seconds, 10)
		}
		if function := logGroup.GetInvokeFunction(); function != nil {
			err := flattenYandexFunctionTriggerInvokeWithRetry(d, function)
			if err != nil {
				return err
			}
		}
		err := d.Set(triggerTypeLogGroup, []map[string]interface{}{lg})
		if err != nil {
			return err
		}
	} else if logging := trig.GetRule().GetLogging(); logging != nil {

		resourceIDs := &schema.Set{F: schema.HashString}
		for _, id := range logging.ResourceId {
			resourceIDs.Add(id)
		}
		resourceTypes := &schema.Set{F: schema.HashString}
		for _, t := range logging.ResourceType {
			resourceTypes.Add(t)
		}
		levels := &schema.Set{F: schema.HashString}
		for _, level := range logging.Levels {
			if l, ok := levelEnumToName[level]; ok {
				levels.Add(l)
			}
		}

		lg := map[string]interface{}{
			"group_id":       logging.LogGroupId,
			"resource_ids":   resourceIDs,
			"resource_types": resourceTypes,
			"levels":         levels,
		}
		if batch := logging.GetBatchSettings(); batch != nil {
			lg["batch_size"] = strconv.FormatInt(batch.Size, 10)
			lg["batch_cutoff"] = strconv.FormatInt(batch.Cutoff.Seconds, 10)
		}
		if function := logging.GetInvokeFunction(); function != nil {
			err := flattenYandexFunctionTriggerInvokeWithRetry(d, function)
			if err != nil {
				return err
			}
		}
		err := d.Set(triggerTypeLogging, []map[string]interface{}{lg})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceYandexFunctionTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutRead))
	defer cancel()

	req := triggers.GetTriggerRequest{
		TriggerId: d.Id(),
	}

	trig, err := config.sdk.Serverless().Triggers().Trigger().Get(ctx, &req)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Yandex Cloud Functions Trigger %q", d.Id()))
	}

	return flattenYandexFunctionTrigger(d, trig)
}

func resourceYandexFunctionTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := config.ContextWithTimeout(d.Timeout(schema.TimeoutDelete))
	defer cancel()

	req := triggers.DeleteTriggerRequest{
		TriggerId: d.Id(),
	}

	op, err := config.sdk.Serverless().Triggers().Trigger().Delete(ctx, &req)
	err = waitOperation(ctx, config, op, err)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Yandex Cloud Functions Trigger %q", d.Id()))
	}

	return nil
}

func expandRetrySettings(d *schema.ResourceData) (*triggers.RetrySettings, error) {
	settings := &triggers.RetrySettings{}
	var err error
	present := false

	if _, ok := d.GetOk("function.0.retry_attempts"); ok {
		present = true
		settings.RetryAttempts, err = strconv.ParseInt(d.Get("function.0.retry_attempts").(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot define function.retry_attempts for Yandex Cloud Functions Trigger: %s", err)
		}
	}

	if _, ok := d.GetOk("function.0.retry_interval"); ok {
		present = true
		retryInterval, err := strconv.ParseInt(d.Get("function.0.retry_interval").(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot define function.retry_interval for Yandex Cloud Functions Trigger: %s", err)
		}
		settings.Interval = &duration.Duration{Seconds: retryInterval}
	}

	if !present {
		return nil, nil
	}

	return settings, nil
}

func expandDLQSettings(d *schema.ResourceData) (*triggers.PutQueueMessage, error) {
	if _, ok := d.GetOk("dlq"); !ok {
		return nil, nil
	}

	settings := &triggers.PutQueueMessage{
		QueueId:          d.Get("dlq.0.queue_id").(string),
		ServiceAccountId: d.Get("dlq.0.service_account_id").(string),
	}

	return settings, nil
}

func expandBatchSettings(d *schema.ResourceData, prefix string) (settings *triggers.BatchSettings, err error) {
	if _, ok := d.GetOk(prefix + ".batch_size"); !ok {
		return nil, nil
	}

	settings = &triggers.BatchSettings{}
	settings.Size, err = strconv.ParseInt(d.Get(prefix+".batch_size").(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Cannot define "+prefix+".batch_size for Yandex Cloud Functions Trigger: %s", err)
	}

	batchCutoff, err := strconv.ParseInt(d.Get(prefix+".batch_cutoff").(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Cannot define "+prefix+"batch_cutoff for Yandex Cloud Functions Trigger: %s", err)
	}
	settings.Cutoff = &duration.Duration{Seconds: batchCutoff}

	return settings, nil
}

func functionTriggerConflictingTypes(triggerType string) []string {
	res := make([]string, 0, len(functionTriggerTypesList)-1)
	for _, tType := range functionTriggerTypesList {
		if tType != triggerType {
			res = append(res, tType)
		}
	}
	return res
}
