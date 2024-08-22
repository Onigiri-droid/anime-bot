// episodeTracker.go
package main

import (
    "encoding/json"
    "os"
)

func saveEpisodeTracker() error {
    data, err := json.Marshal(episodeTracker)
    if err != nil {
        return err
    }
    return os.WriteFile("episodeTracker.json", data, 0644)
}

func loadEpisodeTracker() error {
    data, err := os.ReadFile("episodeTracker.json")
    if err != nil {
        return err
    }
    return json.Unmarshal(data, &episodeTracker)
}
