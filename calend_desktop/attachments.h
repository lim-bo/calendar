#ifndef ATTACHMENTS_H
#define ATTACHMENTS_H


#include "attachments_comp.h"
#include "client.h"

#include <QDialog>
#include <QUrl>
#include <QDragEnterEvent>
#include <QMimeData>
#include <QFileDialog>
#include <QFileInfo>

#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonArray>
namespace Ui {
class attachments;
}

class attachments : public QDialog
{
    Q_OBJECT

public:
    explicit attachments(QString eventID, QString userUID, Client* client, QWidget *parent = nullptr);
    ~attachments();
protected:
    void dragEnterEvent(QDragEnterEvent *event) override;
    void dropEvent(QDropEvent *event) override;

private slots:
    void on_pushButton_clicked();

private:
    Ui::attachments *ui;

    QString eventID;
    QString userUID;
    Client* client;
    QString currentFilePath;
    bool uploadFile(const QString &filePath);
    void loadAttachments();

};

#endif // ATTACHMENTS_H
