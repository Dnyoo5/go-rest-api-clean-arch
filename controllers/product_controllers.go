package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/api-GO/models"
	"example.com/api-GO/utils"
	"github.com/go-chi/chi/v5"
)

type ProductController struct {
	DB *sql.DB
}

		func kirimEmailNotifikasi(namaBarang string) {
			fmt.Printf("Mengirim email notifikasi: %s\n", namaBarang)
			go time.Sleep(5 * time.Second)
			fmt.Printf("✅ Email untuk %s BERHASIL terkirim!\n", namaBarang)
		}

		// GetAllProducts godoc
// @Summary      Menampilkan semua produk
// @Description  Mengambil daftar lengkap produk dari database
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Product
// @Router       /products [get]
		func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")	

		
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			rows, err := c.DB.QueryContext(ctx, "SELECT id, nama, harga, stok FROM products")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()
			
			var katalog []models.Product
			for rows.Next() {
				var p models.Product
				err := rows.Scan(&p.Id, &p.Nama, &p.Harga, &p.Stok)
				if err != nil {
					log.Print("Error Scan", err)
					continue
				}
				katalog = append(katalog, p)
			}

			if err = rows.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(katalog)
		}
		// CreateProduct godoc
// @Summary      Tambah produk baru
// @Description  Menambahkan data produk ke database dengan input JSON
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Data Produk"
// @Success      201      {object}  map[string]string
// @Failure      400      {string}  string "Bad Request"
// @Router       /products [post]
		func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var p models.Product
			err := json.NewDecoder(r.Body).Decode(&p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if errors := utils.ValidateStruct(p); len(errors) > 0 {
				// Kalau ada error, kirim JSON error detail
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"message": "Validasi gagal",
					"errors":  errors,
				})
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			_, err = c.DB.ExecContext(ctx, "INSERT INTO products (nama, harga, stok) VALUES (?, ?, ?)", p.Nama, p.Harga, p.Stok)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			go kirimEmailNotifikasi(p.Nama)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Product created successfully"})
		}
// @Summary      Edit produk baru
// @Description  Memperbarui data produk ke database dengan input JSON
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Data Produk"
// @Success      201      {object}  map[string]string
// @Failure      400      {string}  string "Bad Request"
// @Router       /products/{id} [put]
		func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var p models.Product

			id := chi.URLParam(r, "id") // Pakai Chi lebih ringkas
			
			if id == "" {
				http.Error(w, "ID produk diperlukan", http.StatusBadRequest)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if errors := utils.ValidateStruct(p); len(errors) > 0 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errors)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			_, err = c.DB.ExecContext(ctx, "UPDATE products SET nama=?, harga=?, stok=? WHERE id=?", p.Nama, p.Harga, p.Stok, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
		}
// @Summary      Delete produk baru
// @Description  Menghapus data produk dari database berdasarkan ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Data Produk"
// @Success      201      {object}  map[string]string
// @Failure      400      {string}  string "Bad Request"
// @Router       /products/{id} [delete]
		func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			id := chi.URLParam(r, "id") // Pakai Chi lebih ringkas

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			_, err := c.DB.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			json.NewEncoder(w).Encode(map[string]string{"message": "Sukses hapus produk"})
		}
	