package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ericolvr/maintenance-v2/internal/dto"
	"github.com/ericolvr/maintenance-v2/internal/service"
)

// TestListComparison compara os resultados dos métodos List e ListWithDetails
func TestListComparison(ticketService service.TicketService) {
	ctx := context.Background()
	limit := 5
	offset := 0

	fmt.Println("🔍 Testando comparação entre List() e ListWithDetails()...")
	fmt.Println(strings.Repeat("=", 60))

	// Chamar método original List
	fmt.Println("📋 Executando List()...")
	listResults, listTotal, err := ticketService.List(ctx, limit, offset)
	if err != nil {
		log.Fatalf("Erro no List(): %v", err)
	}

	// Chamar novo método ListWithDetails
	fmt.Println("🚀 Executando ListWithDetails()...")
	detailsResults, detailsTotal, err := ticketService.ListWithDetails(ctx, limit, offset)
	if err != nil {
		log.Fatalf("Erro no ListWithDetails(): %v", err)
	}

	// Comparar totais
	fmt.Printf("📊 Total List(): %d\n", listTotal)
	fmt.Printf("📊 Total ListWithDetails(): %d\n", detailsTotal)
	
	if listTotal != detailsTotal {
		log.Fatalf("❌ ERRO: Totais diferentes! List: %d, ListWithDetails: %d", listTotal, detailsTotal)
	}
	fmt.Println("✅ Totais são iguais")

	// Comparar número de resultados
	if len(listResults) != len(detailsResults) {
		log.Fatalf("❌ ERRO: Número de resultados diferentes! List: %d, ListWithDetails: %d", 
			len(listResults), len(detailsResults))
	}
	fmt.Printf("✅ Número de resultados iguais: %d\n", len(listResults))

	// Comparar cada ticket individualmente
	fmt.Println("\n🔍 Comparando tickets individualmente...")
	
	allEqual := true
	for i := 0; i < len(listResults); i++ {
		listTicket := listResults[i]
		detailsTicket := detailsResults[i]

		fmt.Printf("\n--- Ticket %d (ID: %d) ---\n", i+1, listTicket.ID)
		
		// Comparar campos básicos
		if !compareBasicFields(listTicket, detailsTicket) {
			allEqual = false
			continue
		}

		// Comparar custos
		if !compareCosts(listTicket.Costs, detailsTicket.Costs) {
			allEqual = false
			continue
		}

		// Comparar total de custos
		if listTicket.TotalCost != detailsTicket.TotalCost {
			fmt.Printf("❌ TotalCost diferente: List=%.2f, Details=%.2f\n", 
				listTicket.TotalCost, detailsTicket.TotalCost)
			allEqual = false
			continue
		}

		fmt.Println("✅ Ticket idêntico")
	}

	// Resultado final
	fmt.Println("\n" + strings.Repeat("=", 60))
	if allEqual {
		fmt.Println("🎉 SUCESSO: Todos os dados são EXATAMENTE IGUAIS!")
		fmt.Println("✅ ListWithDetails retorna os mesmos dados que List()")
		
		// Mostrar diferença de performance estimada
		estimatedQueriesOld := 1 + (len(listResults) * 4) // 1 count + 4 queries por ticket
		estimatedQueriesNew := 1 + 1 + len(detailsResults) // 1 count + 1 join + 1 cost query por ticket
		
		fmt.Printf("📈 Performance estimada:\n")
		fmt.Printf("   List(): ~%d queries\n", estimatedQueriesOld)
		fmt.Printf("   ListWithDetails(): ~%d queries\n", estimatedQueriesNew)
		fmt.Printf("   Melhoria: %.1f%% menos queries\n", 
			float64(estimatedQueriesOld-estimatedQueriesNew)/float64(estimatedQueriesOld)*100)
	} else {
		fmt.Println("❌ FALHA: Dados são diferentes!")
		log.Fatal("Os métodos não retornam dados idênticos")
	}
}

func compareBasicFields(list, details dto.TicketResponse) bool {
	equal := true
	
	if list.ID != details.ID {
		fmt.Printf("❌ ID diferente: List=%d, Details=%d\n", list.ID, details.ID)
		equal = false
	}
	
	if list.Number != details.Number {
		fmt.Printf("❌ Number diferente: List=%s, Details=%s\n", list.Number, details.Number)
		equal = false
	}
	
	if list.Status != details.Status {
		fmt.Printf("❌ Status diferente: List=%d, Details=%d\n", list.Status, details.Status)
		equal = false
	}
	
	if list.BranchName != details.BranchName {
		fmt.Printf("❌ BranchName diferente: List=%s, Details=%s\n", list.BranchName, details.BranchName)
		equal = false
	}
	
	// Comparar provider name (pode ser nil)
	if !reflect.DeepEqual(list.ProviderName, details.ProviderName) {
		fmt.Printf("❌ ProviderName diferente: List=%v, Details=%v\n", list.ProviderName, details.ProviderName)
		equal = false
	}
	
	// Comparar distance (pode ser nil)
	if !reflect.DeepEqual(list.Distance, details.Distance) {
		fmt.Printf("❌ Distance diferente: List=%v, Details=%v\n", list.Distance, details.Distance)
		equal = false
	}
	
	return equal
}

func compareCosts(listCosts, detailsCosts []dto.SolutionItemResponse) bool {
	if len(listCosts) != len(detailsCosts) {
		fmt.Printf("❌ Número de custos diferente: List=%d, Details=%d\n", 
			len(listCosts), len(detailsCosts))
		return false
	}
	
	for i, listCost := range listCosts {
		detailsCost := detailsCosts[i]
		
		if listCost.ProblemName != detailsCost.ProblemName ||
		   listCost.SolutionName != detailsCost.SolutionName ||
		   listCost.UnitPrice != detailsCost.UnitPrice ||
		   listCost.Subtotal != detailsCost.Subtotal {
			
			fmt.Printf("❌ Custo %d diferente:\n", i)
			fmt.Printf("   List: %+v\n", listCost)
			fmt.Printf("   Details: %+v\n", detailsCost)
			return false
		}
	}
	
	return true
}

func main() {
	fmt.Println("⚠️  Este é um arquivo de teste para comparar List() vs ListWithDetails()")
	fmt.Println("Para executar o teste, você precisa:")
	fmt.Println("1. Configurar a conexão com o banco de dados")
	fmt.Println("2. Inicializar o ticketService com todas as dependências")
	fmt.Println("3. Chamar TestListComparison(ticketService)")
	fmt.Println("\nEste arquivo demonstra como os dados devem ser comparados.")
}
