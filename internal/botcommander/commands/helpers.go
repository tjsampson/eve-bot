package commands

import (
	"regexp"
	"strings"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// ExtractArtifactsDefinition extracts the ArtifactDefinitions from the CommandOptions
func ExtractArtifactsDefinition(defType string, opts CommandOptions) eveapimodels.ArtifactDefinitions {
	if val, ok := opts[defType]; ok {
		if artifactDefs, ok := val.(eveapimodels.ArtifactDefinitions); ok {
			return artifactDefs
		}
		return nil

	}
	return nil
}

// ExtractBoolOpt extracts a bool key/val from the opts
func ExtractBoolOpt(defType string, opts CommandOptions) bool {
	if val, ok := opts[defType]; ok {
		if forceDepVal, ok := val.(bool); ok {
			return forceDepVal
		}
		return false
	}
	return false
}

// ExtractStringOpt extracts a string key/val from the options
func ExtractStringOpt(defType string, opts CommandOptions) string {
	if val, ok := opts[defType]; ok {
		if envVal, ok := val.(string); ok {
			return envVal
		}
		return ""
	}
	return ""
}

// ExtractStringListOpt extracts a string slice key value from the options
func ExtractStringListOpt(defType string, opts CommandOptions) eveapimodels.StringList {
	if val, ok := opts[defType]; ok {
		if nsVal, ok := val.(string); ok {
			return eveapimodels.StringList{nsVal}
		}
		return nil
	}
	return nil
}

// CleanUrls cleans the incoming URLs
// this iterates the incoming command and removes an encoding slack adds to URLs
func CleanUrls(input string) string {
	matcher := regexp.MustCompile(`<[a-zA-Z]+:\/\/[a-zA-Z._\-:\d\/|?=]+>`)
	matchIndexes := matcher.FindAllStringIndex(input, -1)
	matchCount := len(matchIndexes)

	if matchCount == 0 {
		return input
	}

	cleanPart := input[0:matchIndexes[0][0]]
	for i, v := range matchIndexes {
		if i > 0 {
			previousMatchLastIndex := matchIndexes[i-1][1]
			currentMatchFirstIndex := matchIndexes[i][0]
			middleMatch := input[previousMatchLastIndex:currentMatchFirstIndex]
			cleanPart = cleanPart + middleMatch
		}

		matchedVal := input[v[0]:v[1]]
		cleanVal := ""

		if strings.Contains(matchedVal, "|") {
			vals := strings.Split(matchedVal, "|")
			cleanVal = vals[1][:len(vals[1])-len(">")]
		} else {
			cleanVal = strings.ReplaceAll(matchedVal, "<", "")
			cleanVal = strings.ReplaceAll(cleanVal, ">", "")
		}

		cleanPart = cleanPart + cleanVal
	}
	return cleanPart + input[matchIndexes[matchCount-1][1]:]
}

func hydrateMetadataMap(keyvals []string) params.MetadataMap {
	log.Logger.Debug("hydrateMetadataMap", zap.Strings("keyvals", keyvals))
	result := make(params.MetadataMap, 0)
	if len(keyvals) == 0 {
		return nil
	}
	for _, s := range keyvals {
		log.Logger.Debug("hydrateMetadataMap iteration", zap.String("s", s))
		if strings.Contains(s, "=") {
			log.Logger.Debug("hydrateMetadataMap iteration = ", zap.String("s", s))
			argKV := strings.Split(s, "=")
			log.Logger.Debug("hydrateMetadataMap iteration = argKV", zap.Strings("argKV", argKV))
			key := CleanUrls(argKV[0])
			log.Logger.Debug("hydrateMetadataMap iteration = argKV key", zap.String("key", key))
			value := CleanUrls(strings.Join(argKV[1:], ","))
			log.Logger.Debug("hydrateMetadataMap iteration = argKV value", zap.String("value", value))
			result[key] = value
		}
	}
	return result
}
