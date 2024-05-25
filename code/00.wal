import "fs";
import "path";
import { time } from "system";

struct DirectoryReader {
  let directoryPath: str;

  fn mount(directoryPath: str) {
    this.directoryPath = directoryPath;
  }

  fn readRecentFiles() {
    let allFiles: []str = fs.readDir(this.directoryPath);
    let recentFiles: []str = [];

    foreach file in allFiles {
      let fullPath: str = path.join(this.directoryPath, file);
      let fileInfo: FileInfo = fs.stat(fullPath);
      if this.isFileRecent(fileInfo.creationTime) {
        recentFiles.push(fullPath);
      }
    }

    foreach file in recentFiles {
      println(file, fs.stat(file).creationTime);
    }
  }

  fn isFileRecent(creationTime: Time): bool {
    let twentyFourHoursAgo: Time = time.now() - time.hours(24);
    creationTime > twentyFourHoursAgo;
  }
}

fn main() {
  const directory: str = "/path/to/directory";
  const reader = DirectoryReader::new();
  reader.mount(directory);
  reader.readRecentFiles();
}

main();