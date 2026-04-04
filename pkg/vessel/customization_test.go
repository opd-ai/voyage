package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestGetCustomizationVocab(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		vocab := GetCustomizationVocab(genre)
		if vocab == nil {
			t.Errorf("nil vocab for genre %s", genre)
			continue
		}
		// Test that all fields are populated
		if vocab.ScreenTitle == "" {
			t.Errorf("empty ScreenTitle for genre %s", genre)
		}
		if vocab.SubtitleText == "" {
			t.Errorf("empty SubtitleText for genre %s", genre)
		}
		if vocab.UpgradePointsText == "" {
			t.Errorf("empty UpgradePointsText for genre %s", genre)
		}
		if vocab.ConfirmText == "" {
			t.Errorf("empty ConfirmText for genre %s", genre)
		}
		if vocab.CancelText == "" {
			t.Errorf("empty CancelText for genre %s", genre)
		}
		if vocab.ModuleSectionTitle == "" {
			t.Errorf("empty ModuleSectionTitle for genre %s", genre)
		}
		if vocab.ModuleTierLabel == "" {
			t.Errorf("empty ModuleTierLabel for genre %s", genre)
		}
		if vocab.ModuleEffectLabel == "" {
			t.Errorf("empty ModuleEffectLabel for genre %s", genre)
		}
		if vocab.ModuleUpgradeAction == "" {
			t.Errorf("empty ModuleUpgradeAction for genre %s", genre)
		}
		if vocab.ModuleDowngradeAction == "" {
			t.Errorf("empty ModuleDowngradeAction for genre %s", genre)
		}
		if vocab.LoadoutSectionTitle == "" {
			t.Errorf("empty LoadoutSectionTitle for genre %s", genre)
		}
		if vocab.LoadoutSelectText == "" {
			t.Errorf("empty LoadoutSelectText for genre %s", genre)
		}
		if vocab.VariantSectionTitle == "" {
			t.Errorf("empty VariantSectionTitle for genre %s", genre)
		}
		if vocab.VariantSelectText == "" {
			t.Errorf("empty VariantSelectText for genre %s", genre)
		}
		if vocab.InsigniaSectionTitle == "" {
			t.Errorf("empty InsigniaSectionTitle for genre %s", genre)
		}
		if vocab.InsigniaSelectText == "" {
			t.Errorf("empty InsigniaSelectText for genre %s", genre)
		}
		if vocab.InsuranceSectionTitle == "" {
			t.Errorf("empty InsuranceSectionTitle for genre %s", genre)
		}
		if vocab.InsurancePurchaseText == "" {
			t.Errorf("empty InsurancePurchaseText for genre %s", genre)
		}
		if vocab.SummaryTitle == "" {
			t.Errorf("empty SummaryTitle for genre %s", genre)
		}
		if vocab.TotalCostLabel == "" {
			t.Errorf("empty TotalCostLabel for genre %s", genre)
		}
		if vocab.TotalCapacityLabel == "" {
			t.Errorf("empty TotalCapacityLabel for genre %s", genre)
		}
		if vocab.TotalSpeedLabel == "" {
			t.Errorf("empty TotalSpeedLabel for genre %s", genre)
		}
		if vocab.TotalDefenseLabel == "" {
			t.Errorf("empty TotalDefenseLabel for genre %s", genre)
		}
	}
}

func TestGetModuleUpgradeVocab(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		vocab := GetModuleUpgradeVocab(genre)
		if vocab == nil {
			t.Errorf("nil vocab for genre %s", genre)
			continue
		}
		// Test that all fields are populated
		if vocab.UpgradeText == "" {
			t.Errorf("empty UpgradeText for genre %s", genre)
		}
		if vocab.DowngradeText == "" {
			t.Errorf("empty DowngradeText for genre %s", genre)
		}
		if vocab.RepairText == "" {
			t.Errorf("empty RepairText for genre %s", genre)
		}
		if vocab.SpecializeText == "" {
			t.Errorf("empty SpecializeText for genre %s", genre)
		}
		if vocab.MaxedOutText == "" {
			t.Errorf("empty MaxedOutText for genre %s", genre)
		}
		if vocab.InsufficientText == "" {
			t.Errorf("empty InsufficientText for genre %s", genre)
		}
		if vocab.SuccessText == "" {
			t.Errorf("empty SuccessText for genre %s", genre)
		}
		if vocab.FailureText == "" {
			t.Errorf("empty FailureText for genre %s", genre)
		}
		if vocab.CurrencyName == "" {
			t.Errorf("empty CurrencyName for genre %s", genre)
		}
		if vocab.MaterialsName == "" {
			t.Errorf("empty MaterialsName for genre %s", genre)
		}
		if vocab.SpeedBoostText == "" {
			t.Errorf("empty SpeedBoostText for genre %s", genre)
		}
		if vocab.CargoBoostText == "" {
			t.Errorf("empty CargoBoostText for genre %s", genre)
		}
		if vocab.DefenseBoostText == "" {
			t.Errorf("empty DefenseBoostText for genre %s", genre)
		}
		if vocab.HealingBoostText == "" {
			t.Errorf("empty HealingBoostText for genre %s", genre)
		}
		if vocab.NavigationBoostText == "" {
			t.Errorf("empty NavigationBoostText for genre %s", genre)
		}
	}
}

func TestCustomizationVocab_InvalidGenre(t *testing.T) {
	// Test that invalid genre falls back to fantasy
	vocab := GetCustomizationVocab("invalid")
	fantasyVocab := GetCustomizationVocab(engine.GenreFantasy)

	if vocab.ScreenTitle != fantasyVocab.ScreenTitle {
		t.Error("expected invalid genre to fall back to fantasy")
	}
}

func TestModuleUpgradeVocab_InvalidGenre(t *testing.T) {
	// Test that invalid genre falls back to fantasy
	vocab := GetModuleUpgradeVocab("invalid")
	fantasyVocab := GetModuleUpgradeVocab(engine.GenreFantasy)

	if vocab.UpgradeText != fantasyVocab.UpgradeText {
		t.Error("expected invalid genre to fall back to fantasy")
	}
}

func TestAllCustomizationVocabs(t *testing.T) {
	vocabs := AllCustomizationVocabs()
	if len(vocabs) != 5 {
		t.Errorf("expected 5 vocab sets, got %d", len(vocabs))
	}
}

func TestAllModuleUpgradeVocabs(t *testing.T) {
	vocabs := AllModuleUpgradeVocabs()
	if len(vocabs) != 5 {
		t.Errorf("expected 5 vocab sets, got %d", len(vocabs))
	}
}

func TestCustomizationVocab_GenreUniqueness(t *testing.T) {
	// Verify that each genre has unique vocabulary
	genres := engine.AllGenres()
	for i, g1 := range genres {
		for j, g2 := range genres {
			if i >= j {
				continue
			}
			v1 := GetCustomizationVocab(g1)
			v2 := GetCustomizationVocab(g2)

			// Screen titles should be different across genres
			if v1.ScreenTitle == v2.ScreenTitle {
				t.Errorf("genres %s and %s have same ScreenTitle", g1, g2)
			}
		}
	}
}

func TestModuleUpgradeVocab_GenreUniqueness(t *testing.T) {
	// Verify that each genre has unique vocabulary
	genres := engine.AllGenres()
	for i, g1 := range genres {
		for j, g2 := range genres {
			if i >= j {
				continue
			}
			v1 := GetModuleUpgradeVocab(g1)
			v2 := GetModuleUpgradeVocab(g2)

			// Currency names should be different across genres
			if v1.CurrencyName == v2.CurrencyName {
				t.Errorf("genres %s and %s have same CurrencyName", g1, g2)
			}
		}
	}
}
