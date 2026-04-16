<?php

$repoName = "warframe-public-export-plus";
$commitHash = "b0dedf4a0b666cf95d570e281ef63dc3de3a8d21";
$url = "https://github.com/calamity-inc/$repoName/archive/$commitHash.zip";
$inputPath = __DIR__ . "/$repoName.zip";
$destPath = __DIR__ . "/";

if (downloadAndExtractZip($url, $inputPath, $destPath, $repoName)) {
	echo "Download and extraction successful.\n";
} else {
	echo "Something went wrong.\n";
}

function downloadAndExtractZip(
	string $url,
	string $zipPath,
	string $extractTo,
	string $newRoot,
): bool {
	if (!downloadFile($url, $zipPath)) {
		echo "Download failed.\n";
		return false;
	}

	if (!is_dir($extractTo)) {
		if (!mkdir($extractTo, 0777, true)) {
			echo "Failed to create extract directory.\n";
			return false;
		}
	}

	$zip = new ZipArchive();
	if ($zip->open($zipPath) !== true) {
		echo "Failed to open ZIP archive.\n";
		return false;
	}

	if (!extractZip($zip, $extractTo, $newRoot)) {
		echo "Extraction failed.\n";
		$zip->close();
		return false;
	}

	$zip->close();
	return true;
}

function extractZip(ZipArchive $zip, string $extractTo, string $newRoot): bool
{
	$roots = [];
	for ($i = 0; $i < $zip->numFiles; $i++) {
		$name = str_replace("\\", "/", $zip->getNameIndex($i));
		if ($name === "") {
			continue;
		}

		$parts = explode("/", $name);
		if ($parts[0] !== "") {
			$roots[$parts[0]] = true;
		}
	}

	$flatten = null;
	if (count($roots) === 1) {
		$flatten = array_key_first($roots);
	}

	for ($i = 0; $i < $zip->numFiles; $i++) {
		$name = str_replace("\\", "/", $zip->getNameIndex($i));
		if ($name === "") {
			continue;
		}

		$relPath = $name;

		if ($flatten !== null && str_starts_with($name, $flatten . "/")) {
			$relPath = substr($name, strlen($flatten) + 1);
		}

		$relPath = ltrim($relPath, "/");
		$relPath = str_replace("./", "", $relPath);

		if (
			$relPath === "" ||
			str_contains($relPath, "../") ||
			str_contains($relPath, "..\\")
		) {
			continue;
		}

		$target = rtrim($extractTo, "/\\") . "/" . $newRoot . "/" . $relPath;
		if (str_ends_with($name, "/")) {
			if (!is_dir($target)) {
				mkdir($target, 0777, true);
			}
			continue;
		}

		$dir = dirname($target);
		if (!is_dir($dir)) {
			mkdir($dir, 0777, true);
		}

		$in = $zip->getStream($name);
		if (!$in) {
			continue;
		}

		$out = fopen($target, "wb");
		if (!$out) {
			fclose($in);
			continue;
		}

		while (!feof($in)) {
			fwrite($out, fread($in, 8192));
		}

		fclose($in);
		fclose($out);
	}

	return true;
}

function downloadFile(string $url, string $dest): bool
{
	$in = fopen($url, "rb");
	if (!$in) {
		return false;
	}

	$out = fopen($dest, "wb");
	if (!$out) {
		fclose($in);
		return false;
	}

	while (!feof($in)) {
		fwrite($out, fread($in, 8192));
	}

	fclose($in);
	fclose($out);

	return true;
}
