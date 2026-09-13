/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Eoghain2708/dreamreadr/cmd"
	"github.com/Eoghain2708/dreamreadr/internal/ai"
	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/Eoghain2708/dreamreadr/internal/service"
	"github.com/Eoghain2708/dreamreadr/internal/storage"
	_ "modernc.org/sqlite"
)

func main() {
	var analyser dream.DreamAnalyser
	var embedder dream.DreamEmbedder

	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := storage.NewSQLRepository(db)
	analyser, err = ai.NewLocalLLMAnalyser(http.DefaultClient, "http://127.0.0.1:8080")
	if err != nil {
		log.Fatalf("cannot create analyser")
	}

	embedder, err = ai.NewLocalLLMEmbedder("http://127.0.0.1:8081", "nomic-ai/nomic-embed-text-v1.5-GGUF:Q4_K_M")
	if err != nil {
		log.Fatalf("cannot create embedder")
	}

	ds := service.NewDreamService(repo, repo, repo, repo, analyser, embedder)
	cmd.SetDreamService(ds)
	cmd.Execute()
}
