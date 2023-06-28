package product

import (
	"context"
	"net/http"

	//"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"errors"

	"github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/emicklei/go-restful/v3"
	//"log"
)

type EgwProductHttpHandler struct {
	productSvc ports.EgwProductUsecase
}

func NewEgwProductHandler(productSvc ports.EgwProductUsecase, wsCont *restful.Container) *EgwProductHttpHandler {
	httpHandler := &EgwProductHttpHandler{
		productSvc: productSvc,
	}

	ws := new(restful.WebService)

	ws.Path("/product").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/insert").To(httpHandler.InsertProduct))
	ws.Route(ws.PUT("/update/{id}").To(httpHandler.UpdateProduct))
	//ws.Route(ws.DELETE("/delete/{id}").To(httpHandler.DeleteProduct))

	wsCont.Add(ws)

	return httpHandler
}

// Performs insert product
func (e *EgwProductHttpHandler) InsertProduct(req *restful.Request, resp *restful.Response) {
	var reqData InsertRequestData
	req.ReadEntity(&reqData)

	var egwProduct *EgwProductModel = &EgwProductModel{}
	egwProduct.Name = reqData.Name
	egwProduct.ShortDescription = reqData.ShortDescription
	egwProduct.Description = reqData.Description
	egwProduct.Price = reqData.Price

	err := e.insertProduct(req.Request.Context(), egwProduct)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("error insert product"))
		return
	}

	// send product back
	respData := InsertResponseData{Product: *egwProduct}

	resp.WriteAsJson(respData)
}

func (e *EgwProductHttpHandler) insertProduct(ctx context.Context, egwProduct *EgwProductModel) error {
	err := e.productSvc.InsertProduct(ctx, egwProduct.ToDomain())
	if err != nil {
		return err
	}
	// retrieve their data from the DB to populate it (e.g. ID)
	productData, err := e.productSvc.FindByID(ctx, egwProduct.ID)
	if err != nil {
		return err
	}
	egwProduct.FromDomain(productData)
	return nil
}

func (e *EgwProductHttpHandler) UpdateProduct(req *restful.Request, resp *restful.Response) {

	var a UpdateRequestData
	req.ReadEntity(&a)

	productID := req.PathParameter("id")

	existingProduct, err := e.productSvc.FindByID(req.Request.Context(), productID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("failed to retrieve product"))
		return
	}

	existingProduct.ShortDescription = a.ShortDescription
	existingProduct.Description = a.Description
	existingProduct.Price = a.Price

	var egwProduct *EgwProductModel = &EgwProductModel{}
	egwProduct.ID = existingProduct.ID
	egwProduct.Name = existingProduct.Name
	egwProduct.ShortDescription = existingProduct.ShortDescription
	egwProduct.Description = existingProduct.Description
	egwProduct.Price = existingProduct.Price
	egwProduct.CreatedAt = existingProduct.CreatedAt
	egwProduct.UpdatedAt = existingProduct.UpdatedAt

	err = e.updateProduct(req.Request.Context(), egwProduct)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("failed to update product"))
		return
	}

	// return updated product as data
	respData := UpdateResponseData{Product: *egwProduct}
	resp.WriteAsJson(respData)
}

func (e *EgwProductHttpHandler) updateProduct(ctx context.Context, egwProduct *EgwProductModel) error {
	err := e.productSvc.Update(ctx, egwProduct.ToDomain())
	if err != nil {
		return err
	}
	// Optionally retrieve the updated product data from the DB to populate it further, if needed
	updatedProductData, err := e.productSvc.FindByID(ctx, egwProduct.ID)
	if err != nil {
		return err
	}
	egwProduct.FromDomain(updatedProductData)
	return nil
}

func (e *EgwProductHttpHandler) DeleteProduct(req *restful.Request, resp *restful.Response) {
	// Retrieve the product ID from the URL path
	productID := req.PathParameter("id")

	// Delete the product from the database
	err := e.deleteProduct(req.Request.Context(), productID)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, errors.New("failed to delete product"))
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Product deleted successfully"))
}

func (e *EgwProductHttpHandler) deleteProduct(ctx context.Context, productID string) error {
	err := e.productSvc.Delete(ctx, productID)
	if err != nil {
		return err
	}
	return nil
}
