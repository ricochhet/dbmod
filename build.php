<?php

$sourceDir =
	__DIR__ . DIRECTORY_SEPARATOR . "warframe-public-export-plus-senpai";
$destBase = __DIR__ . DIRECTORY_SEPARATOR . "assets";

$sourceDir = realpath($sourceDir);
$iterator = new RecursiveIteratorIterator(
	new RecursiveDirectoryIterator(
		$sourceDir,
		RecursiveDirectoryIterator::SKIP_DOTS,
	),
);

foreach ($iterator as $file) {
	if ($file->getExtension() !== "json") {
		continue;
	}
	$fullPath = $file->getRealPath();
	$relPath = substr($fullPath, strlen($sourceDir) + 1);
	$destPath = $destBase . DIRECTORY_SEPARATOR . $relPath;
	$destDir = dirname($destPath);
	if (!is_dir($destDir)) {
		mkdir($destDir, 0777, true);
	}
	copy($fullPath, $destPath);
}
