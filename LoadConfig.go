// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srcom

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

/*

Ver.0.0.0
Ver.1.0.0 LoginShowroom()の戻り値 status を err に変更する。
Ver.-.-.- exsrapi.go から分離する。
Ver.1.1.0 yaml.Unmarshal() を yaml.UnmarshalStrict() に変更する。

*/

// 設定ファイルを読み込む
//
//	以下の記事を参考にさせていただきました。
//		【Go初学】設定ファイル、環境変数から設定情報を取得する
//			https://note.com/artefactnote/n/n8c22d1ac4b86
func LoadConfig(filePath string, config interface{}) (err error) {

	var content []byte

	if strings.Contains(filepath.Base(filePath), ".enc.") {
		if os.Getenv("SOPS_AGE_KEY_FILE") == "" {
			return fmt.Errorf("SOPS_AGE_KEY_FILE is empty")
		}

		content, err = exec.Command("sops", "--decrypt", filePath).Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return fmt.Errorf("sops --decrypt failed: %w: %s", err, strings.TrimSpace(string(exitErr.Stderr)))
			}
			return fmt.Errorf("sops --decrypt failed: %w", err)
		}
	} else {
		content, err = os.ReadFile(filePath)
		if err != nil {
			err = fmt.Errorf("os.ReadFile: %w", err)
			return err
		}
	}

	content = []byte(os.ExpandEnv(string(content)))
	//	log.Printf("content=%s\n", content)

	if err := yaml.UnmarshalStrict(content, config); err != nil {
		err = fmt.Errorf("yaml.UnmarshalStrict(): %w", err)
		return err
	}

	//	log.Printf("\n")
	//	log.Printf("%+v\n", config)
	//	log.Printf("\n")

	return nil
}
