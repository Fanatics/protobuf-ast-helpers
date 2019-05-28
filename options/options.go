// Package options provides convenience methods for reading node's options.
package options

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/richardwilkes/toolbox/errs"
)

// OptionReader embeds an extension registry
type OptionReader struct {
	init     bool
	Registry *dynamic.ExtensionRegistry
}

// NewOptionReader uses a default extension registry if the passed one is nil
func NewOptionReader(registry *dynamic.ExtensionRegistry) *OptionReader {
	if registry == nil {
		registry = dynamic.NewExtensionRegistryWithDefaults()
	}

	o := &OptionReader{
		Registry: registry,
	}

	return o
}

// ReadOptionByID attempts to load the value into the "out" param if it is the correct type,
// if the found type and the out type do not agree an error is returned
func (o *OptionReader) ReadOptionByID(node desc.Descriptor, tagID int, out interface{}) error {
	val, err := o.GetOptionByID(node, tagID)
	if err != nil {
		return errs.Wrap(err)
	}

	return o.readValue(val, out)
}

// ReadOptionByName attempts to load the value into the "out" param if it is the correct type,
// if the found type and the out type do not agree an error is returned
func (o *OptionReader) ReadOptionByName(node desc.Descriptor, tagName string, out interface{}) error {
	val, err := o.GetOptionByName(node, tagName)
	if err != nil {
		return errs.Wrap(err)
	}

	return o.readValue(val, out)
}

// GetOptionByID attempts to find an option on the node by tagID
func (o *OptionReader) GetOptionByID(node desc.Descriptor, tagID int) (interface{}, error) {
	opts, err := o.getOptions(node)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	val, err := opts.TryGetFieldByNumber(tagID)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return val, nil
}

// GetOptionByName attempts to find an option on the node by tagName
func (o *OptionReader) GetOptionByName(node desc.Descriptor, tagName string) (interface{}, error) {
	opts, err := o.getOptions(node)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	val, err := opts.TryGetFieldByName(tagName)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return val, nil
}

// -------
// helpers

func (o *OptionReader) getOptions(node desc.Descriptor) (*dynamic.Message, error) {
	if node == nil {
		return nil, nil
	}

	if !o.init {
		o.Registry.AddExtensionsFromFileRecursively(node.GetFile())
		o.init = true
	}

	opts := node.GetOptions()
	if opts == nil {
		return nil, nil
	}

	msg, err := dynamic.AsDynamicMessageWithExtensionRegistry(opts, o.Registry)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return msg, nil
}

func (o *OptionReader) readValue(in interface{}, out interface{}) error {
	if in == nil {
		return nil
	}
	if out == nil {
		return nil
	}

	switch target := out.(type) {
	case *string:
		data, ok := in.(*string)
		if !ok {
			return errs.Newf("value was %T, not %T", out, data)
		}
		if data != nil {
			*target = *data
		}
	case *bool:
		data, ok := in.(*bool)
		if !ok {
			return errs.Newf("value was %T, not %T", out, data)
		}
		if data != nil {
			*target = *data
		}
	case *int16:
		data, ok := in.(*int16)
		if !ok {
			return errs.Newf("value was %T, not %T", out, data)
		}
		if data != nil {
			*target = *data
		}
	case *int32:
		data, ok := in.(*int32)
		if !ok {
			return errs.Newf("value was %T, not %T", out, data)
		}
		if data != nil {
			*target = *data
		}
	case *int64:
		data, ok := in.(*int64)
		if !ok {
			return errs.Newf("value was %T, not %T", out, data)
		}
		if data != nil {
			*target = *data
		}
	default:
		return errs.Newf("helper could not read %T, please use `GetOption...` directly", target)
	}

	return nil
}
