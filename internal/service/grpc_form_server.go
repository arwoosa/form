package service

import (
	"context"

	"github.com/arwoosa/form-service/gen/pb/common"
	pb "github.com/arwoosa/form-service/gen/pb/form"
	"github.com/arwoosa/form-service/internal/models"
	"github.com/arwoosa/vulpes/ezgrpc"
	"github.com/arwoosa/vulpes/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCFormServer implements the FormService gRPC interface
type GRPCFormServer struct {
	pb.UnimplementedFormServiceServer
	templateService *FormTemplateService
	formService     *FormService
}

// NewGRPCFormServer creates a new gRPC form server
func NewGRPCFormServer(templateService *FormTemplateService, formService *FormService) *GRPCFormServer {
	return &GRPCFormServer{
		templateService: templateService,
		formService:     formService,
	}
}

// CreateFormTemplate creates a new form template
func (s *GRPCFormServer) CreateFormTemplate(ctx context.Context, req *pb.CreateFormTemplateRequest) (*pb.CreateFormTemplateResponse, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Convert request to service input
	input := &models.CreateFormTemplateInput{
		Name:        req.Name,
		Description: req.Description,
		MerchantID:  user.Merchant,
		CreatedBy:   user.ID,
	}

	// Convert schema if provided
	if req.Schema != nil {
		input.Schema = req.Schema.AsMap()
	}

	// Convert UI schema if provided
	if req.Uischema != nil {
		input.UISchema = req.Uischema.AsMap()
	}

	template, err := s.templateService.CreateTemplate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	pbTemplate, err := s.convertFormTemplateToProto(template)
	if err != nil {
		log.Error("Failed to convert template to protobuf", log.Err(err))
		return nil, err
	}

	return &pb.CreateFormTemplateResponse{
		Template: pbTemplate,
	}, nil
}

// ListFormTemplates lists form templates with pagination
func (s *GRPCFormServer) ListFormTemplates(ctx context.Context, req *pb.ListFormTemplatesRequest) (*pb.ListFormTemplatesResponse, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Convert request to service options
	options := &models.FormTemplateQueryOptions{
		MerchantID: user.Merchant,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
	}

	templates, totalCount, err := s.templateService.ListTemplates(ctx, options)
	if err != nil {
		return nil, err
	}

	// Convert templates to protobuf
	pbTemplates := make([]*pb.FormTemplate, len(templates))
	for i, template := range templates {
		pbTemplate, err := s.convertFormTemplateToProto(template)
		if err != nil {
			log.Error("Failed to convert template to protobuf", log.Err(err))
			return nil, err
		}
		pbTemplates[i] = pbTemplate
	}

	// Calculate pagination
	totalPages := (totalCount + int64(options.PageSize) - 1) / int64(options.PageSize)

	return &pb.ListFormTemplatesResponse{
		Templates: pbTemplates,
		Pagination: &common.Pagination{
			Page:       int32(options.Page),
			PageSize:   int32(options.PageSize),
			TotalCount: int32(totalCount),
			TotalPages: int32(totalPages),
		},
	}, nil
}

// GetFormTemplate gets a form template by ID
func (s *GRPCFormServer) GetFormTemplate(ctx context.Context, req *common.ID) (*pb.FormTemplate, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	templateID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	template, err := s.templateService.GetTemplate(ctx, templateID, user.Merchant)
	if err != nil {
		return nil, err
	}

	return s.convertFormTemplateToProto(template)
}

// UpdateFormTemplate updates a form template
func (s *GRPCFormServer) UpdateFormTemplate(ctx context.Context, req *pb.UpdateFormTemplateRequest) (*pb.FormTemplate, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	templateID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	// Convert request to service input
	input := &models.UpdateFormTemplateInput{
		ID:          templateID,
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   user.ID,
		MerchantID:  user.Merchant,
	}

	// Convert schema if provided
	if req.Schema != nil {
		input.Schema = req.Schema.AsMap()
	}

	// Convert UI schema if provided
	if req.Uischema != nil {
		input.UISchema = req.Uischema.AsMap()
	}

	template, err := s.templateService.UpdateTemplate(ctx, input)
	if err != nil {
		return nil, err
	}

	return s.convertFormTemplateToProto(template)
}

// DeleteFormTemplate deletes a form template
func (s *GRPCFormServer) DeleteFormTemplate(ctx context.Context, req *common.ID) (*emptypb.Empty, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	templateID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	err = s.templateService.DeleteTemplate(ctx, templateID, user.Merchant)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DuplicateFormTemplate duplicates a form template
func (s *GRPCFormServer) DuplicateFormTemplate(ctx context.Context, req *pb.DuplicateFormTemplateRequest) (*pb.DuplicateFormTemplateResponse, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	sourceID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	input := &models.DuplicateFormTemplateInput{
		SourceID:   sourceID,
		Name:       req.Name,
		CreatedBy:  user.ID,
		MerchantID: user.Merchant,
	}

	template, err := s.templateService.DuplicateTemplate(ctx, input)
	if err != nil {
		return nil, err
	}

	pbTemplate, err := s.convertFormTemplateToProto(template)
	if err != nil {
		return nil, err
	}

	return &pb.DuplicateFormTemplateResponse{
		Template: pbTemplate,
	}, nil
}

// CreateForm creates a new form
func (s *GRPCFormServer) CreateForm(ctx context.Context, req *pb.CreateFormRequest) (*pb.CreateFormResponse, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Convert request to service input
	input := &models.CreateFormInput{
		Name:        req.Name,
		Description: req.Description,
		MerchantID:  user.Merchant,
		CreatedBy:   user.ID,
	}

	// Convert optional template ID
	if req.TemplateId != "" {
		templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		input.TemplateID = &templateID
	}

	// Convert optional event ID
	if req.EventId != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventId)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		input.EventID = &eventID
	}

	// Convert schema if provided
	if req.Schema != nil {
		input.Schema = req.Schema.AsMap()
	}

	// Convert UI schema if provided
	if req.Uischema != nil {
		input.UISchema = req.Uischema.AsMap()
	}

	form, err := s.formService.CreateForm(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	pbForm, err := s.convertFormToProto(form)
	if err != nil {
		return nil, err
	}

	return &pb.CreateFormResponse{
		Form: pbForm,
	}, nil
}

// ListForms lists forms with pagination and filters
func (s *GRPCFormServer) ListForms(ctx context.Context, req *pb.ListFormsRequest) (*pb.ListFormsResponse, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	// Convert request to service options
	options := &models.FormQueryOptions{
		MerchantID: user.Merchant,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
	}

	// Add optional filters
	if req.EventId != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventId)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		options.EventID = &eventID
	}

	if req.TemplateId != "" {
		templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		options.TemplateID = &templateID
	}

	forms, totalCount, err := s.formService.ListForms(ctx, options)
	if err != nil {
		return nil, err
	}

	// Convert forms to protobuf
	pbForms := make([]*pb.Form, len(forms))
	for i, form := range forms {
		pbForm, err := s.convertFormToProto(form)
		if err != nil {
			return nil, err
		}
		pbForms[i] = pbForm
	}

	// Calculate pagination
	totalPages := (totalCount + int64(options.PageSize) - 1) / int64(options.PageSize)

	return &pb.ListFormsResponse{
		Forms: pbForms,
		Pagination: &common.Pagination{
			Page:       int32(options.Page),
			PageSize:   int32(options.PageSize),
			TotalCount: int32(totalCount),
			TotalPages: int32(totalPages),
		},
	}, nil
}

// GetForm gets a form by ID
func (s *GRPCFormServer) GetForm(ctx context.Context, req *common.ID) (*pb.Form, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	formID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	form, err := s.formService.GetForm(ctx, formID, user.Merchant)
	if err != nil {
		return nil, err
	}

	return s.convertFormToProto(form)
}

// UpdateForm updates a form
func (s *GRPCFormServer) UpdateForm(ctx context.Context, req *pb.UpdateFormRequest) (*pb.Form, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	formID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	// Convert request to service input
	input := &models.UpdateFormInput{
		ID:          formID,
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   user.ID,
		MerchantID:  user.Merchant,
	}

	// Convert optional template ID
	if req.TemplateId != "" {
		templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		input.TemplateID = &templateID
	}

	// Convert optional event ID
	if req.EventId != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventId)
		if err != nil {
			return nil, ErrInvalidObjectID
		}
		input.EventID = &eventID
	}

	// Convert schema if provided
	if req.Schema != nil {
		input.Schema = req.Schema.AsMap()
	}

	// Convert UI schema if provided
	if req.Uischema != nil {
		input.UISchema = req.Uischema.AsMap()
	}

	form, err := s.formService.UpdateForm(ctx, input)
	if err != nil {
		return nil, err
	}

	return s.convertFormToProto(form)
}

// DeleteForm deletes a form
func (s *GRPCFormServer) DeleteForm(ctx context.Context, req *common.ID) (*emptypb.Empty, error) {
	user, err := ezgrpc.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	formID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	err = s.formService.DeleteForm(ctx, formID, user.Merchant)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// convertFormTemplateToProto converts a form template model to protobuf
func (s *GRPCFormServer) convertFormTemplateToProto(template *models.FormTemplate) (*pb.FormTemplate, error) {
	var schemaStruct *structpb.Struct
	var uiSchemaStruct *structpb.Struct
	var err error

	if template.Schema != nil {
		schemaMap := s.convertMongoDataToMap(template.Schema)
		if schemaMap != nil {
			schemaStruct, err = structpb.NewStruct(schemaMap)
			if err != nil {
				return nil, err
			}
		}
	}

	if template.UISchema != nil {
		uiSchemaMap := s.convertMongoDataToMap(template.UISchema)
		if uiSchemaMap != nil {
			uiSchemaStruct, err = structpb.NewStruct(uiSchemaMap)
			if err != nil {
				return nil, err
			}
		}
	}

	return &pb.FormTemplate{
		Id:          template.ID.Hex(),
		Name:        template.Name,
		MerchantId:  template.MerchantID,
		Description: template.Description,
		Schema:      schemaStruct,
		Uischema:    uiSchemaStruct,
		CreatedAt:   timestamppb.New(template.GetCreatedAt()),
		CreatedBy:   template.CreatedBy,
		UpdatedAt:   timestamppb.New(template.GetUpdatedAt()),
		UpdatedBy:   template.UpdatedBy,
	}, nil
}

// convertFormToProto converts a form model to protobuf
func (s *GRPCFormServer) convertFormToProto(form *models.Form) (*pb.Form, error) {
	var schemaStruct *structpb.Struct
	var uiSchemaStruct *structpb.Struct
	var err error

	if form.Schema != nil {
		schemaMap := s.convertMongoDataToMap(form.Schema)
		if schemaMap != nil {
			schemaStruct, err = structpb.NewStruct(schemaMap)
			if err != nil {
				return nil, err
			}
		}
	}

	if form.UISchema != nil {
		uiSchemaMap := s.convertMongoDataToMap(form.UISchema)
		if uiSchemaMap != nil {
			uiSchemaStruct, err = structpb.NewStruct(uiSchemaMap)
			if err != nil {
				return nil, err
			}
		}
	}

	pbForm := &pb.Form{
		Id:          form.ID.Hex(),
		Name:        form.Name,
		MerchantId:  form.MerchantID,
		Description: form.Description,
		Schema:      schemaStruct,
		Uischema:    uiSchemaStruct,
		CreatedAt:   timestamppb.New(form.GetCreatedAt()),
		CreatedBy:   form.CreatedBy,
		UpdatedAt:   timestamppb.New(form.GetUpdatedAt()),
		UpdatedBy:   form.UpdatedBy,
	}

	// Add optional fields
	if form.TemplateID != nil {
		pbForm.TemplateId = form.TemplateID.Hex()
	}

	if form.EventID != nil {
		pbForm.EventId = form.EventID.Hex()
	}

	return pbForm, nil
}

// convertMongoDataToMap converts MongoDB primitive types to map[string]interface{}
func (s *GRPCFormServer) convertMongoDataToMap(data interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case primitive.D:
		result := make(map[string]interface{})
		for _, elem := range v {
			result[elem.Key] = s.convertValue(elem.Value)
		}
		return result
	case map[string]interface{}:
		// Already correct type, but need to recursively process values
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = s.convertValue(value)
		}
		return result
	case primitive.M:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = s.convertValue(value)
		}
		return result
	default:
		// If not a map type, try to convert to single value map
		if converted := s.convertValue(v); converted != nil {
			return map[string]interface{}{"value": converted}
		}
		return nil
	}
}

// convertValue converts single MongoDB values, handling nested primitive.D, primitive.A etc
func (s *GRPCFormServer) convertValue(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case primitive.D:
		result := make(map[string]interface{})
		for _, elem := range v {
			result[elem.Key] = s.convertValue(elem.Value)
		}
		return result
	case primitive.A:
		result := make([]interface{}, len(v))
		for i, elem := range v {
			result[i] = s.convertValue(elem)
		}
		return result
	case primitive.M:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = s.convertValue(value)
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = s.convertValue(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, elem := range v {
			result[i] = s.convertValue(elem)
		}
		return result
	default:
		return data
	}
}
