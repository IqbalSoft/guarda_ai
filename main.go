package main

import (
	"fmt"
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
)

type TrayUI struct {
	ZaiH, Zai5HTitle, Zai5H, Zai5HR, ZaiWkTitle, ZaiWk, ZaiWkR, ZaiMoTitle, ZaiMo, ZaiMoR, ZaiSep *systray.MenuItem
	ClaudeH, ClaudeQ, ClaudeSep                                                                   *systray.MenuItem
	GPT_H, GPT_Q, GPT_Sep                                                                         *systray.MenuItem
	DeepH, DeepQ, DeepSep                                                                         *systray.MenuItem
	GemH, GemQ, GemSep                                                                            *systray.MenuItem
	KimiH, KimiQ, KimiSep                                                                         *systray.MenuItem
}

var ui TrayUI

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	iconData, err := os.ReadFile("icon.png")
	if err == nil {
		systray.SetIcon(iconData)
	}

	systray.SetTitle("AI: Syncing...")
	systray.SetTooltip("GuardaAI")

	// --- MAIN MENU (Z.AI) ---
	ui.ZaiH = systray.AddMenuItem("=== Z.ai GLM ===", "")
	ui.ZaiH.Disable()
	ui.Zai5HTitle = systray.AddMenuItem("5-Hour Quota", "")
	ui.Zai5HTitle.Disable()
	ui.Zai5H = systray.AddMenuItem("[Loading...]", "")
	ui.Zai5H.Disable()
	ui.Zai5HR = systray.AddMenuItem("Reset Time: Loading...", "")
	ui.Zai5HR.Disable()
	ui.ZaiWkTitle = systray.AddMenuItem("Weekly Quota", "")
	ui.ZaiWkTitle.Disable()
	ui.ZaiWk = systray.AddMenuItem("[Loading...]", "")
	ui.ZaiWk.Disable()
	ui.ZaiWkR = systray.AddMenuItem("Reset: Loading...", "")
	ui.ZaiWkR.Disable()
	ui.ZaiMoTitle = systray.AddMenuItem("Monthly Quota", "")
	ui.ZaiMoTitle.Disable()
	ui.ZaiMo = systray.AddMenuItem("[Loading...]", "")
	ui.ZaiMo.Disable()
	ui.ZaiMoR = systray.AddMenuItem("Reset: Loading...", "")
	ui.ZaiMoR.Disable()
	ui.ZaiSep = systray.AddMenuItem("────────────────────────", "")
	ui.ZaiSep.Disable()

	// --- SUBMENU: SEE MORE ---
	mMoreAIs := systray.AddMenuItem("See More AIs...", "View other tools")
	ui.ClaudeH, ui.ClaudeQ, ui.ClaudeSep = createSubGroup(mMoreAIs, "=== Claude ===")
	ui.GPT_H, ui.GPT_Q, ui.GPT_Sep = createSubGroup(mMoreAIs, "=== ChatGPT ===")
	ui.DeepH, ui.DeepQ, ui.DeepSep = createSubGroup(mMoreAIs, "=== Deepseek ===")
	ui.GemH, ui.GemQ, ui.GemSep = createSubGroup(mMoreAIs, "=== Gemini ===")
	ui.KimiH, ui.KimiQ, ui.KimiSep = createSubGroup(mMoreAIs, "=== Kimi ===")
	systray.AddMenuItem("────────────────────────", "").Disable()

	// --- SETTINGS SUBMENU ---
	mSettings := systray.AddMenuItem("Settings", "Configure plugins")
	mPluginsMenu := mSettings.AddSubMenuItem("Plugins", "")

	mZaiVisibility, mZaiAPIKey := createPluginSettingsRow(mPluginsMenu, "Z.ai GLM", true)
	mClaudeVisibility, _ := createPluginSettingsRow(mPluginsMenu, "Claude", false)
	mChatGPTVisibility, _ := createPluginSettingsRow(mPluginsMenu, "ChatGPT", false)
	mDeepseekVisibility, _ := createPluginSettingsRow(mPluginsMenu, "Deepseek", false)
	mGeminiVisibility, _ := createPluginSettingsRow(mPluginsMenu, "Gemini", false)
	mKimiVisibility, _ := createPluginSettingsRow(mPluginsMenu, "Kimi", false)

	go func() {
		cfg := loadConfig()
		updatePluginSettingsState(mZaiVisibility, mZaiAPIKey, "Z.ai GLM", cfg.Zai.Enabled, cfg.Zai.APIKey, true)
		updatePluginSettingsState(mClaudeVisibility, nil, "Claude", cfg.Claude.Enabled, "", false)
		updatePluginSettingsState(mChatGPTVisibility, nil, "ChatGPT", cfg.ChatGPT.Enabled, "", false)
		updatePluginSettingsState(mDeepseekVisibility, nil, "Deepseek", cfg.Deepseek.Enabled, "", false)
		updatePluginSettingsState(mGeminiVisibility, nil, "Gemini", cfg.Gemini.Enabled, "", false)
		updatePluginSettingsState(mKimiVisibility, nil, "Kimi", cfg.Kimi.Enabled, "", false)

		for {
			select {
			case <-mZaiVisibility.ClickedCh:
				currentCfg := loadConfig()
				currentCfg.Zai.Enabled = !currentCfg.Zai.Enabled
				saveConfig(currentCfg)
				updatePluginSettingsState(mZaiVisibility, mZaiAPIKey, "Z.ai GLM", currentCfg.Zai.Enabled, currentCfg.Zai.APIKey, true)
				refreshAllData()

			case <-mClaudeVisibility.ClickedCh:
				currentCfg := loadConfig()
				currentCfg.Claude.Enabled = !currentCfg.Claude.Enabled
				saveConfig(currentCfg)
				updatePluginSettingsState(mClaudeVisibility, nil, "Claude", currentCfg.Claude.Enabled, "", false)
				refreshAllData()

			case <-mChatGPTVisibility.ClickedCh:
				currentCfg := loadConfig()
				currentCfg.ChatGPT.Enabled = !currentCfg.ChatGPT.Enabled
				saveConfig(currentCfg)
				updatePluginSettingsState(mChatGPTVisibility, nil, "ChatGPT", currentCfg.ChatGPT.Enabled, "", false)
				refreshAllData()

			case <-mDeepseekVisibility.ClickedCh:
				currentCfg := loadConfig()
				currentCfg.Deepseek.Enabled = !currentCfg.Deepseek.Enabled
				saveConfig(currentCfg)
				updatePluginSettingsState(mDeepseekVisibility, nil, "Deepseek", currentCfg.Deepseek.Enabled, "", false)
				refreshAllData()

			case <-mGeminiVisibility.ClickedCh:
				currentCfg := loadConfig()
				currentCfg.Gemini.Enabled = !currentCfg.Gemini.Enabled
				saveConfig(currentCfg)
				updatePluginSettingsState(mGeminiVisibility, nil, "Gemini", currentCfg.Gemini.Enabled, "", false)
				refreshAllData()

			case <-mKimiVisibility.ClickedCh:
				currentCfg := loadConfig()
				currentCfg.Kimi.Enabled = !currentCfg.Kimi.Enabled
				saveConfig(currentCfg)
				updatePluginSettingsState(mKimiVisibility, nil, "Kimi", currentCfg.Kimi.Enabled, "", false)
				refreshAllData()

			case <-mZaiAPIKey.ClickedCh:
				currentCfg := loadConfig()
				newKey, err := zenity.Entry(
					"Enter API Key for Z.ai GLM:",
					zenity.Title("Settings - Z.ai GLM"),
					zenity.EntryText(currentCfg.Zai.APIKey),
				)
				if err == nil {
					currentCfg.Zai.APIKey = newKey
					saveConfig(currentCfg)
					updatePluginSettingsState(mZaiVisibility, mZaiAPIKey, "Z.ai GLM", currentCfg.Zai.Enabled, currentCfg.Zai.APIKey, true)
					refreshAllData()
				}
			}
		}
	}()

	systray.AddMenuItem("────────────────────────", "").Disable()

	// --- CONTROLS ---
	mRefresh := systray.AddMenuItem("Refresh Data", "Force refresh")
	mQuit := systray.AddMenuItem("Quit", "Exit application")

	applyVisibility()
	refreshAllData()

	// Click Handlers
	go func() {
		for {
			select {
			case <-mRefresh.ClickedCh:
				refreshAllData()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// Auto-Refresh
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			refreshAllData()
		}
	}()
}

func createPluginSettingsRow(parent *systray.MenuItem, title string, hasApiKeyAction bool) (*systray.MenuItem, *systray.MenuItem) {
	toggle := parent.AddSubMenuItem("[ ] "+title, "Show in tray")

	var action *systray.MenuItem
	if hasApiKeyAction {
		action = parent.AddSubMenuItem("Edit API Key", "Update your token")
	}

	return toggle, action
}

func updatePluginSettingsState(toggle, action *systray.MenuItem, title string, enabled bool, apiKey string, hasApiKeyAction bool) {
	if enabled {
		toggle.SetTitle("[x] " + title)
	} else {
		toggle.SetTitle("[ ] " + title)
	}

	if action != nil && hasApiKeyAction {
		if enabled {
			action.Show()
			action.Enable()
		} else {
			action.Hide()
			action.Disable()
		}
	}
}

func createSubGroup(parent *systray.MenuItem, title string) (*systray.MenuItem, *systray.MenuItem, *systray.MenuItem) {
	h := parent.AddSubMenuItem(title, "")
	h.Disable()
	q := parent.AddSubMenuItem("Quota: [Loading...]", "")
	q.Disable()
	sep := parent.AddSubMenuItem("────────────────────────", "")
	sep.Disable()
	return h, q, sep
}

func applyVisibility() {
	cfg := loadConfig()

	if ui.ZaiH != nil {
		if cfg.Zai.Enabled {
			ui.ZaiH.Show()
			ui.Zai5HTitle.Show()
			ui.Zai5H.Show()
			ui.Zai5HR.Show()
			ui.ZaiWkTitle.Show()
			ui.ZaiWk.Show()
			ui.ZaiWkR.Show()
			ui.ZaiMoTitle.Show()
			ui.ZaiMo.Show()
			ui.ZaiMoR.Show()
			ui.ZaiSep.Show()
		} else {
			ui.ZaiH.Hide()
			ui.Zai5HTitle.Hide()
			ui.Zai5H.Hide()
			ui.Zai5HR.Hide()
			ui.ZaiWkTitle.Hide()
			ui.ZaiWk.Hide()
			ui.ZaiWkR.Hide()
			ui.ZaiMoTitle.Hide()
			ui.ZaiMo.Hide()
			ui.ZaiMoR.Hide()
			ui.ZaiSep.Hide()
		}
	}

	toggleGroup(cfg.Claude.Enabled, ui.ClaudeH, ui.ClaudeQ, ui.ClaudeSep)
	toggleGroup(cfg.ChatGPT.Enabled, ui.GPT_H, ui.GPT_Q, ui.GPT_Sep)
	toggleGroup(cfg.Deepseek.Enabled, ui.DeepH, ui.DeepQ, ui.DeepSep)
	toggleGroup(cfg.Gemini.Enabled, ui.GemH, ui.GemQ, ui.GemSep)
	toggleGroup(cfg.Kimi.Enabled, ui.KimiH, ui.KimiQ, ui.KimiSep)
}

func toggleGroup(enabled bool, h, q, sep *systray.MenuItem) {
	if h == nil || q == nil || sep == nil {
		return
	}
	if enabled {
		h.Show()
		q.Show()
		sep.Show()
	} else {
		h.Hide()
		q.Hide()
		sep.Hide()
	}
}

func refreshAllData() {
	systray.SetTitle("AI: Syncing...")
	cfg := loadConfig()
	applyVisibility()

	fetchZaiUsage(ui.Zai5H, ui.Zai5HR, ui.ZaiWk, ui.ZaiWkR, ui.ZaiMo, ui.ZaiMoR, cfg.Zai)
	fetchGenericAI("Claude", ui.ClaudeQ, cfg.Claude)
	fetchGenericAI("ChatGPT", ui.GPT_Q, cfg.ChatGPT)
	fetchGenericAI("Deepseek", ui.DeepQ, cfg.Deepseek)
	fetchGenericAI("Gemini", ui.GemQ, cfg.Gemini)
	fetchGenericAI("Kimi", ui.KimiQ, cfg.Kimi)

	systray.SetTitle("GuardaAI")
}

func onExit() {
	fmt.Println("Shutting down GuardaAI...")
}
