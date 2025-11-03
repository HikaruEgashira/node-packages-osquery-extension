-- Example queries for node_packages table

-- 1. List all packages
SELECT * FROM node_packages;

-- 2. Count packages by manager
SELECT manager, COUNT(*) as count
FROM node_packages
GROUP BY manager;

-- 3. Find specific package (e.g., React)
SELECT * FROM node_packages
WHERE name LIKE '%react%';

-- 4. List npm packages only
SELECT name, version, cache_path
FROM node_packages
WHERE manager = 'npm';

-- 5. Show unique packages with version count
SELECT name, COUNT(DISTINCT version) as version_count,
       GROUP_CONCAT(DISTINCT manager) as managers
FROM node_packages
GROUP BY name
ORDER BY version_count DESC;

-- 6. Find packages with specific version
SELECT name, manager, cache_path
FROM node_packages
WHERE version LIKE '1.0.%';

-- 7. Count total unique packages
SELECT COUNT(DISTINCT name) as unique_packages
FROM node_packages;

-- 8. Find packages in multiple managers
SELECT name, COUNT(DISTINCT manager) as manager_count,
       GROUP_CONCAT(DISTINCT manager) as managers
FROM node_packages
GROUP BY name
HAVING manager_count > 1;

-- 9. List Deno packages
SELECT name, version
FROM node_packages
WHERE manager = 'deno';

-- 10. Search by cache path
SELECT name, version, manager
FROM node_packages
WHERE cache_path LIKE '%/.npm/%';

-- 11. Find latest version of each package per manager
SELECT name, manager, MAX(version) as latest_version
FROM node_packages
GROUP BY name, manager
ORDER BY name;

-- 12. Packages with @scope (scoped packages)
SELECT name, version, manager
FROM node_packages
WHERE name LIKE '@%';

-- 13. Statistics per manager
SELECT
    manager,
    COUNT(*) as total_packages,
    COUNT(DISTINCT name) as unique_packages
FROM node_packages
GROUP BY manager;

-- 14. Find TypeScript related packages
SELECT name, version, manager
FROM node_packages
WHERE name LIKE '%typescript%' OR name LIKE '%ts-%';

-- 15. Most cached packages (across all managers)
SELECT name, COUNT(*) as cache_count,
       GROUP_CONCAT(DISTINCT version) as versions
FROM node_packages
GROUP BY name
ORDER BY cache_count DESC
LIMIT 10;
