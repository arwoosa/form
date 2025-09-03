package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/arwoosa/form-service/gen/pb/common"
	pb "github.com/arwoosa/form-service/gen/pb/form"
	"github.com/arwoosa/form-service/internal/models"
	"github.com/arwoosa/vulpes/log"
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

// Helper functions for conversion
func (s *GRPCFormServer) modelToProtoTemplate(template *models.FormTemplate) (*pb.FormTemplate, error) {
	schema, err := structpb.NewStruct(template.Schema.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("failed to convert schema: %w", err)
	}

	var uiSchema *structpb.Struct
	if template.UISchema != nil {
		uiSchema, err = structpb.NewStruct(template.UISchema.(map[string]interface{}))
		if err != nil {
			return nil, fmt.Errorf("failed to convert ui_schema: %w", err)
		}
	}

	return &pb.FormTemplate{
		Id:          template.ID.Hex(),
		Name:        template.Name,
		MerchantId:  template.MerchantID,
		Description: template.Description,
		Schema:      schema,
		UiSchema:    uiSchema,
		CreatedAt:   timestamppb.New(template.GetCreatedAt()),
		CreatedBy:   template.CreatedBy,
		UpdatedAt:   timestamppb.New(template.GetUpdatedAt()),
		UpdatedBy:   template.UpdatedBy,
	}, nil
}

func (s *GRPCFormServer) modelToProtoForm(form *models.Form) (*pb.Form, error) {
	schema, err := structpb.NewStruct(form.Schema.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("failed to convert schema: %w", err)
	}

	var uiSchema *structpb.Struct
	if form.UISchema != nil {
		uiSchema, err = structpb.NewStruct(form.UISchema.(map[string]interface{}))
		if err != nil {
			return nil, fmt.Errorf("failed to convert ui_schema: %w", err)
		}
	}

	result := &pb.Form{
		Id:          form.ID.Hex(),
		Name:        form.Name,
		MerchantId:  form.MerchantID,
		Description: form.Description,
		Schema:      schema,
		UiSchema:    uiSchema,
		CreatedAt:   timestamppb.New(form.GetCreatedAt()),
		CreatedBy:   form.CreatedBy,
		UpdatedAt:   timestamppb.New(form.GetUpdatedAt()),
		UpdatedBy:   form.UpdatedBy,
	}

	if form.EventID != nil && !form.EventID.IsZero() {
		result.EventId = form.EventID.Hex()
	}

	if form.TemplateID != nil && !form.TemplateID.IsZero() {
		result.TemplateId = form.TemplateID.Hex()
	}

	return result, nil
}

func (s *GRPCFormServer) createSuccessResponse() *common.BaseResponse {
	return &common.BaseResponse{
		Success: true,
		Error:   nil,
	}
}

func (s *GRPCFormServer) createErrorResponse(err error) *common.BaseResponse {
	return &common.BaseResponse{
		Success: false,
		Error: &common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	}
}

// Form Template Operations
func (s *GRPCFormServer) CreateFormTemplate(ctx context.Context, req *pb.CreateFormTemplateRequest) (*pb.CreateFormTemplateResponse, error) {
	log.Info("CreateFormTemplate called", log.String("name", req.Name))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		log.Error("Failed to get user info", log.Err(err))
		return &pb.CreateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert proto request to model input
	input := &models.CreateFormTemplateInput{
		Name:        req.Name,
		Description: req.Description,
		Schema:      req.Schema.AsMap(),
		UISchema:    req.UiSchema.AsMap(),
		CreatedBy:   userInfo.UserID,
		MerchantID:  userInfo.MerchantID,
	}

	// Create template
	template, err := s.templateService.CreateTemplate(ctx, input)
	if err != nil {
		log.Error("Failed to create template", log.Err(err))
		return &pb.CreateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoTemplate, err := s.modelToProtoTemplate(template)
	if err != nil {
		log.Error("Failed to convert template to proto", log.Err(err))
		return &pb.CreateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.CreateFormTemplateResponse{
		Base:     s.createSuccessResponse(),
		Template: protoTemplate,
	}, nil
}

func (s *GRPCFormServer) GetFormTemplate(ctx context.Context, req *pb.GetFormTemplateRequest) (*pb.GetFormTemplateResponse, error) {
	log.Info("GetFormTemplate called", log.String("template_id", req.TemplateId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.GetFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
	if err != nil {
		return &pb.GetFormTemplateResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Get template
	template, err := s.templateService.GetTemplate(ctx, templateID, userInfo.MerchantID)
	if err != nil {
		return &pb.GetFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoTemplate, err := s.modelToProtoTemplate(template)
	if err != nil {
		return &pb.GetFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.GetFormTemplateResponse{
		Base:     s.createSuccessResponse(),
		Template: protoTemplate,
	}, nil
}

func (s *GRPCFormServer) ListFormTemplates(ctx context.Context, req *pb.ListFormTemplatesRequest) (*pb.ListFormTemplatesResponse, error) {
	log.Info("ListFormTemplates called")

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.ListFormTemplatesResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert request to query options
	options := &models.FormTemplateQueryOptions{
		MerchantID: userInfo.MerchantID,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
	}

	// List templates
	templates, totalCount, err := s.templateService.ListTemplates(ctx, options)
	if err != nil {
		return &pb.ListFormTemplatesResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoTemplates := make([]*pb.FormTemplate, len(templates))
	for i, template := range templates {
		protoTemplate, err := s.modelToProtoTemplate(template)
		if err != nil {
			return &pb.ListFormTemplatesResponse{
				Base: s.createErrorResponse(err),
			}, ToGRPCError(err)
		}
		protoTemplates[i] = protoTemplate
	}

	// Calculate pagination
	totalPages := int32((totalCount + int64(options.PageSize) - 1) / int64(options.PageSize))
	pagination := &common.Pagination{
		Page:       int32(options.Page),
		PageSize:   int32(options.PageSize),
		TotalCount: int32(totalCount),
		TotalPages: totalPages,
	}

	return &pb.ListFormTemplatesResponse{
		Base:       s.createSuccessResponse(),
		Templates:  protoTemplates,
		Pagination: pagination,
	}, nil
}

func (s *GRPCFormServer) UpdateFormTemplate(ctx context.Context, req *pb.UpdateFormTemplateRequest) (*pb.UpdateFormTemplateResponse, error) {
	log.Info("UpdateFormTemplate called", log.String("template_id", req.TemplateId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.UpdateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
	if err != nil {
		return &pb.UpdateFormTemplateResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Convert request to model input
	input := &models.UpdateFormTemplateInput{
		ID:          templateID,
		Name:        req.Name,
		Description: req.Description,
		Schema:      req.Schema.AsMap(),
		UISchema:    req.UiSchema.AsMap(),
		UpdatedBy:   userInfo.UserID,
	}

	// Update template
	template, err := s.templateService.UpdateTemplate(ctx, input)
	if err != nil {
		return &pb.UpdateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoTemplate, err := s.modelToProtoTemplate(template)
	if err != nil {
		return &pb.UpdateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.UpdateFormTemplateResponse{
		Base:     s.createSuccessResponse(),
		Template: protoTemplate,
	}, nil
}

func (s *GRPCFormServer) DeleteFormTemplate(ctx context.Context, req *pb.DeleteFormTemplateRequest) (*pb.DeleteFormTemplateResponse, error) {
	log.Info("DeleteFormTemplate called", log.String("template_id", req.TemplateId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.DeleteFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
	if err != nil {
		return &pb.DeleteFormTemplateResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Delete template
	err = s.templateService.DeleteTemplate(ctx, templateID, userInfo.MerchantID)
	if err != nil {
		return &pb.DeleteFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.DeleteFormTemplateResponse{
		Base: s.createSuccessResponse(),
	}, nil
}

func (s *GRPCFormServer) DuplicateFormTemplate(ctx context.Context, req *pb.DuplicateFormTemplateRequest) (*pb.DuplicateFormTemplateResponse, error) {
	log.Info("DuplicateFormTemplate called", log.String("template_id", req.TemplateId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.DuplicateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
	if err != nil {
		return &pb.DuplicateFormTemplateResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Convert request to model input
	input := &models.DuplicateFormTemplateInput{
		SourceID:   templateID,
		Name:       req.Name,
		CreatedBy:  userInfo.UserID,
		MerchantID: userInfo.MerchantID,
	}

	// Duplicate template
	template, err := s.templateService.DuplicateTemplate(ctx, input)
	if err != nil {
		return &pb.DuplicateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoTemplate, err := s.modelToProtoTemplate(template)
	if err != nil {
		return &pb.DuplicateFormTemplateResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.DuplicateFormTemplateResponse{
		Base:     s.createSuccessResponse(),
		Template: protoTemplate,
	}, nil
}

// Form Operations (continuing in the same pattern)
func (s *GRPCFormServer) CreateForm(ctx context.Context, req *pb.CreateFormRequest) (*pb.CreateFormResponse, error) {
	log.Info("CreateForm called", log.String("name", req.Name))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.CreateFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert request to model input
	input := &models.CreateFormInput{
		Name:        req.Name,
		Description: req.Description,
		Schema:      req.Schema.AsMap(),
		UISchema:    req.UiSchema.AsMap(),
		CreatedBy:   userInfo.UserID,
		MerchantID:  userInfo.MerchantID,
	}

	// Handle optional EventID
	if req.EventId != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventId)
		if err != nil {
			return &pb.CreateFormResponse{
				Base: s.createErrorResponse(ErrInvalidInput),
			}, ToGRPCError(ErrInvalidInput)
		}
		input.EventID = &eventID
	}

	// Handle optional TemplateID
	if req.TemplateId != "" {
		templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
		if err != nil {
			return &pb.CreateFormResponse{
				Base: s.createErrorResponse(ErrInvalidInput),
			}, ToGRPCError(ErrInvalidInput)
		}
		input.TemplateID = &templateID
	}

	// Create form
	form, err := s.formService.CreateForm(ctx, input)
	if err != nil {
		return &pb.CreateFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoForm, err := s.modelToProtoForm(form)
	if err != nil {
		return &pb.CreateFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.CreateFormResponse{
		Base: s.createSuccessResponse(),
		Form: protoForm,
	}, nil
}

func (s *GRPCFormServer) GetForm(ctx context.Context, req *pb.GetFormRequest) (*pb.GetFormResponse, error) {
	log.Info("GetForm called", log.String("form_id", req.FormId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.GetFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	formID, err := primitive.ObjectIDFromHex(req.FormId)
	if err != nil {
		return &pb.GetFormResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Get form
	form, err := s.formService.GetForm(ctx, formID, userInfo.MerchantID)
	if err != nil {
		return &pb.GetFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoForm, err := s.modelToProtoForm(form)
	if err != nil {
		return &pb.GetFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.GetFormResponse{
		Base: s.createSuccessResponse(),
		Form: protoForm,
	}, nil
}

func (s *GRPCFormServer) ListForms(ctx context.Context, req *pb.ListFormsRequest) (*pb.ListFormsResponse, error) {
	log.Info("ListForms called")

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.ListFormsResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert request to query options
	options := &models.FormQueryOptions{
		MerchantID: userInfo.MerchantID,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
	}

	// Handle optional filters
	if req.EventId != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventId)
		if err != nil {
			return &pb.ListFormsResponse{
				Base: s.createErrorResponse(ErrInvalidInput),
			}, ToGRPCError(ErrInvalidInput)
		}
		options.EventID = &eventID
	}

	if req.TemplateId != "" {
		templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
		if err != nil {
			return &pb.ListFormsResponse{
				Base: s.createErrorResponse(ErrInvalidInput),
			}, ToGRPCError(ErrInvalidInput)
		}
		options.TemplateID = &templateID
	}

	// List forms
	forms, totalCount, err := s.formService.ListForms(ctx, options)
	if err != nil {
		return &pb.ListFormsResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoForms := make([]*pb.Form, len(forms))
	for i, form := range forms {
		protoForm, err := s.modelToProtoForm(form)
		if err != nil {
			return &pb.ListFormsResponse{
				Base: s.createErrorResponse(err),
			}, ToGRPCError(err)
		}
		protoForms[i] = protoForm
	}

	// Calculate pagination
	totalPages := int32((totalCount + int64(options.PageSize) - 1) / int64(options.PageSize))
	pagination := &common.Pagination{
		Page:       int32(options.Page),
		PageSize:   int32(options.PageSize),
		TotalCount: int32(totalCount),
		TotalPages: totalPages,
	}

	return &pb.ListFormsResponse{
		Base:       s.createSuccessResponse(),
		Forms:      protoForms,
		Pagination: pagination,
	}, nil
}

func (s *GRPCFormServer) UpdateForm(ctx context.Context, req *pb.UpdateFormRequest) (*pb.UpdateFormResponse, error) {
	log.Info("UpdateForm called", log.String("form_id", req.FormId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.UpdateFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	formID, err := primitive.ObjectIDFromHex(req.FormId)
	if err != nil {
		return &pb.UpdateFormResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Convert request to model input
	input := &models.UpdateFormInput{
		ID:          formID,
		Name:        req.Name,
		Description: req.Description,
		Schema:      req.Schema.AsMap(),
		UISchema:    req.UiSchema.AsMap(),
		UpdatedBy:   userInfo.UserID,
	}

	// Handle optional EventID
	if req.EventId != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventId)
		if err != nil {
			return &pb.UpdateFormResponse{
				Base: s.createErrorResponse(ErrInvalidInput),
			}, ToGRPCError(ErrInvalidInput)
		}
		input.EventID = &eventID
	}

	// Handle optional TemplateID
	if req.TemplateId != "" {
		templateID, err := primitive.ObjectIDFromHex(req.TemplateId)
		if err != nil {
			return &pb.UpdateFormResponse{
				Base: s.createErrorResponse(ErrInvalidInput),
			}, ToGRPCError(ErrInvalidInput)
		}
		input.TemplateID = &templateID
	}

	// Update form
	form, err := s.formService.UpdateForm(ctx, input)
	if err != nil {
		return &pb.UpdateFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert to proto
	protoForm, err := s.modelToProtoForm(form)
	if err != nil {
		return &pb.UpdateFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.UpdateFormResponse{
		Base: s.createSuccessResponse(),
		Form: protoForm,
	}, nil
}

func (s *GRPCFormServer) DeleteForm(ctx context.Context, req *pb.DeleteFormRequest) (*pb.DeleteFormResponse, error) {
	log.Info("DeleteForm called", log.String("form_id", req.FormId))

	// Get user info
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return &pb.DeleteFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	// Convert ID
	formID, err := primitive.ObjectIDFromHex(req.FormId)
	if err != nil {
		return &pb.DeleteFormResponse{
			Base: s.createErrorResponse(ErrInvalidInput),
		}, ToGRPCError(ErrInvalidInput)
	}

	// Delete form
	err = s.formService.DeleteForm(ctx, formID, userInfo.MerchantID)
	if err != nil {
		return &pb.DeleteFormResponse{
			Base: s.createErrorResponse(err),
		}, ToGRPCError(err)
	}

	return &pb.DeleteFormResponse{
		Base: s.createSuccessResponse(),
	}, nil
}
