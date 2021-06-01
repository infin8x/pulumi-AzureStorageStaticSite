import * as AzureStorageStaticSite from "@pulumi/AzureStorageStaticSite";

const page = new AzureStorageStaticSite.StaticPage("page", {
    indexContent: "<html><body><p>Hello world!</p></body></html>",
});

export const storageAccount = page.storageAccount;
export const websiteUrl = page.websiteUrl;
export const cdnUrl = page.cdnUrl;
